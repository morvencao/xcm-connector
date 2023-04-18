package cluster

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/morvencao/xcm-connector/pkg/helpers"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
	operatorhelpers "github.com/openshift/library-go/pkg/operator/v1helpers"

	clusterclient "open-cluster-management.io/api/client/cluster/clientset/versioned"
	clusterinformerv1 "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1"
	clusterlisterv1 "open-cluster-management.io/api/client/cluster/listers/cluster/v1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const ManagedClusterConditionConnected string = "ManagedClusterConditionConnected"

type labelsChange struct {
	changedLabels    map[string]string
	removedLabelKeys []string
}

type clusterController struct {
	clusterClient clusterclient.Interface
	clusterLister clusterlisterv1.ManagedClusterLister
	queue         workqueue.RateLimitingInterface
	xCMAPIServer  string
	accessToken   string
	// TODO: remove this if the label/status change can be extracted from reconcile event
	clusterLabelsChangeMap      map[string]labelsChange
	clusterLabelsChangeMapMutex *sync.RWMutex
	clusterStatusChangeMap      map[string][]metav1.Condition
	clusterStatusChangeMapMutex *sync.RWMutex
}

func NewClusterController(
	clusterClient clusterclient.Interface,
	clusterInformer clusterinformerv1.ManagedClusterInformer,
	xCMAPIServer string,
	accessToken string,
	recorder events.Recorder,
) factory.Controller {
	controllerName := "cluster-controller"
	syncCtx := factory.NewSyncContext(controllerName, recorder)

	c := &clusterController{
		clusterClient:               clusterClient,
		clusterLister:               clusterInformer.Lister(),
		queue:                       syncCtx.Queue(),
		xCMAPIServer:                xCMAPIServer,
		accessToken:                 accessToken,
		clusterLabelsChangeMap:      make(map[string]labelsChange),
		clusterLabelsChangeMapMutex: &sync.RWMutex{},
		clusterStatusChangeMap:      make(map[string][]metav1.Condition),
		clusterStatusChangeMapMutex: &sync.RWMutex{},
	}

	clusterEventHandlerFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cluster, ok := obj.(*clusterv1.ManagedCluster)
			if !ok {
				runtime.HandleError(fmt.Errorf("error to get object: %v", obj))
				return
			}
			c.queue.Add(cluster.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// only need handle label and conditions update
			oldCluster, ok := oldObj.(*clusterv1.ManagedCluster)
			if !ok {
				runtime.HandleError(fmt.Errorf("error to get object: %v", oldObj))
				return
			}
			newCluster, ok := newObj.(*clusterv1.ManagedCluster)
			if !ok {
				runtime.HandleError(fmt.Errorf("error to get object: %v", newObj))
				return
			}
			enqueue := false
			// TODO: ignore build-in labels
			if !reflect.DeepEqual(oldCluster.Labels, newCluster.Labels) {
				c.updateClusterChangeLabels(newCluster.GetName(), oldCluster.Labels, newCluster.Labels)
				enqueue = true
			}
			if !reflect.DeepEqual(oldCluster.Status.Conditions, newCluster.Status.Conditions) {
				c.updateClusterChangeStatus(newCluster.GetName(), oldCluster.Status.Conditions, newCluster.Status.Conditions)
				enqueue = true
			}
			if enqueue {
				c.queue.Add(newCluster.GetName())
			}
		},
		DeleteFunc: func(obj interface{}) {
			switch t := obj.(type) {
			case *clusterv1.ManagedCluster:
				c.queue.Add(t.GetName())
			case cache.DeletedFinalStateUnknown:
				cluster, ok := t.Obj.(*clusterv1.ManagedCluster)
				if !ok {
					runtime.HandleError(fmt.Errorf("error to get object: %v", obj))
					return
				}
				c.queue.Add(cluster.GetName())
			default:
				runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
				return
			}
		},
	}

	_, err := clusterInformer.Informer().AddEventHandler(clusterEventHandlerFuncs)

	if err != nil {
		runtime.HandleError(err)
	}

	return factory.New().
		WithSyncContext(syncCtx).
		// WithInformersQueueKeyFunc(func(obj runtime.Object) string {
		// 	accessor, _ := meta.Accessor(obj)
		// 	return accessor.GetName()
		// }, clusterInformer.Informer()).
		WithBareInformers(clusterInformer.Informer()).
		WithSync(c.sync).
		ToController(controllerName, recorder)
}

func (c *clusterController) sync(ctx context.Context, syncCtx factory.SyncContext) error {
	managedClusterName := syncCtx.QueueKey()
	klog.Infof("Reconciling ManagedCluster %s", managedClusterName)

	cluster, err := c.clusterLister.Get(managedClusterName)
	if errors.IsNotFound(err) {
		// Spoke cluster not found, could have been deleted, do nothing.
		return nil
	}
	if err != nil {
		return err
	}

	// check condition firstly
	cond := meta.FindStatusCondition(cluster.Status.Conditions, ManagedClusterConditionConnected)
	if cond != nil {
		var errList []error

		// update cluster labels
		c.clusterLabelsChangeMapMutex.RLock()
		labelsChanges := c.clusterLabelsChangeMap[managedClusterName]
		c.clusterLabelsChangeMapMutex.RUnlock()

		if len(labelsChanges.changedLabels) > 0 {
			klog.Infof("posting cluster labels: %s", labelsChanges.changedLabels)
			if err := helpers.PostClusterLabels(c.xCMAPIServer, c.accessToken, cluster, labelsChanges.changedLabels); err != nil {
				errList = append(errList, fmt.Errorf("failed to update labels %s for cluster %s: %v", labelsChanges.changedLabels, cluster.GetName(), err))
			} else {
				labelsChanges.changedLabels = make(map[string]string)
			}
			klog.Infof("after posting cluster labels: %s", labelsChanges.changedLabels)
		}
		if len(labelsChanges.removedLabelKeys) > 0 {
			klog.Infof("removing cluster label keys: %s", labelsChanges.removedLabelKeys)
			if err := helpers.RemoveClusterLabels(c.xCMAPIServer, c.accessToken, cluster, labelsChanges.removedLabelKeys); err != nil {
				errList = append(errList, fmt.Errorf("failed to remove label keys %s for cluster %s: %v", labelsChanges.removedLabelKeys, cluster.GetName(), err))
			} else {
				labelsChanges.removedLabelKeys = make([]string, 0)
			}
			klog.Infof("after removing cluster label keys: %s", labelsChanges.removedLabelKeys)
		}

		c.clusterLabelsChangeMapMutex.Lock()
		// reset the labels change
		c.clusterLabelsChangeMap[managedClusterName] = labelsChanges
		c.clusterLabelsChangeMapMutex.Unlock()

		// update cluster status
		c.clusterStatusChangeMapMutex.RLock()
		addContidions := c.clusterStatusChangeMap[managedClusterName]
		c.clusterStatusChangeMapMutex.RUnlock()
		if len(addContidions) > 0 {
			if err := helpers.PostClusterStatus(c.xCMAPIServer, c.accessToken, cluster, addContidions); err != nil {
				errList = append(errList, fmt.Errorf("failed to add status conditions %v for cluster %s: %v", addContidions, cluster.GetName(), err))
			} else {
				addContidions = []metav1.Condition{}
			}
		}

		c.clusterStatusChangeMapMutex.Lock()
		// reset the status change
		c.clusterStatusChangeMap[managedClusterName] = addContidions
		c.clusterStatusChangeMapMutex.Unlock()

		if len(errList) > 0 {
			return operatorhelpers.NewMultiLineAggregate(errList)
		}

		return nil
	}

	// TODO timeout or retry limits
	if err := helpers.CreateCluster(c.xCMAPIServer, c.accessToken, cluster); err != nil {
		return err
	}

	conditionUpdateFn := helpers.UpdateManagedClusterConditionFn(metav1.Condition{
		Type:    ManagedClusterConditionConnected,
		Status:  metav1.ConditionTrue,
		Reason:  "ManagedClusterConnected",
		Message: fmt.Sprintf("Managed cluster %s is connected to xCM.", cluster.Name),
	})

	_, updated, err := helpers.UpdateManagedClusterStatus(ctx, c.clusterClient, cluster.Name, conditionUpdateFn)
	if updated {
		syncCtx.Recorder().Eventf("ManagedClusterConnectedConditionUpdated",
			"update managed cluster %q connected condition to true",
			cluster.Name)
	}

	return err
}

// updateClusterChangeLabels
func (c *clusterController) updateClusterChangeLabels(clusterName string, old, new map[string]string) {
	c.clusterLabelsChangeMapMutex.RLock()
	labelsChanges := c.clusterLabelsChangeMap[clusterName]
	c.clusterLabelsChangeMapMutex.RUnlock()

	for k, v := range new {
		_, oldExists := old[k]
		if !oldExists {
			if labelsChanges.changedLabels == nil {
				labelsChanges.changedLabels = make(map[string]string)
			}
			labelsChanges.changedLabels[k] = v
		}
	}

	for k := range old {
		if _, exists := new[k]; !exists {
			// check if in changedLabels
			if _, existsInChangedLabels := labelsChanges.changedLabels[k]; existsInChangedLabels {
				// remove from changedLabels
				delete(labelsChanges.changedLabels, k)
			}
			if labelsChanges.removedLabelKeys == nil {
				labelsChanges.removedLabelKeys = make([]string, 0)
			}
			labelsChanges.removedLabelKeys = append(labelsChanges.removedLabelKeys, k)
		}
	}

	klog.Infof("update cluster labels: %s", labelsChanges.changedLabels)
	klog.Infof("update cluster labels removed keys: %s", labelsChanges.removedLabelKeys)
	c.clusterLabelsChangeMapMutex.Lock()
	c.clusterLabelsChangeMap[clusterName] = labelsChanges
	c.clusterLabelsChangeMapMutex.Unlock()
}

// updateClusterChangeStatus
func (c *clusterController) updateClusterChangeStatus(clusterName string, old, new []metav1.Condition) {
	c.clusterStatusChangeMapMutex.RLock()
	addContidions := c.clusterStatusChangeMap[clusterName]
	c.clusterStatusChangeMapMutex.RUnlock()

	for _, cond := range new {
		oldCond := meta.FindStatusCondition(old, cond.Type)
		if oldCond != nil && oldCond.Status == cond.Status {
			continue
		}
		meta.SetStatusCondition(&addContidions, cond)
	}

	c.clusterStatusChangeMapMutex.Lock()
	c.clusterStatusChangeMap[clusterName] = addContidions
	c.clusterStatusChangeMapMutex.Unlock()
}
