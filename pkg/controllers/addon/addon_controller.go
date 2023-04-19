package addon

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/morvencao/xcm-connector/pkg/helpers"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	addonclient "open-cluster-management.io/api/client/addon/clientset/versioned"
	addoninformerv1alpha1 "open-cluster-management.io/api/client/addon/informers/externalversions/addon/v1alpha1"
	addonlisterv1alpha1 "open-cluster-management.io/api/client/addon/listers/addon/v1alpha1"
	clusterinformerv1 "open-cluster-management.io/api/client/cluster/informers/externalversions/cluster/v1"
	clusterlisterv1 "open-cluster-management.io/api/client/cluster/listers/cluster/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const ManagedClusterAddonConditionConnected string = "ManagedClusterAddonConditionConnected"

type addonController struct {
	clusterLister clusterlisterv1.ManagedClusterLister
	addonClient   addonclient.Interface
	addonLister   addonlisterv1alpha1.ManagedClusterAddOnLister
	queue         workqueue.RateLimitingInterface
	xCMAPIServer  string
	accessToken   string
	// TODO: remove this if the status change can be extracted from reconcile event
	addonStatusChangeMap      map[string][]metav1.Condition
	addonStatusChangeMapMutex *sync.RWMutex
}

func NewAddonController(
	clusterInformer clusterinformerv1.ManagedClusterInformer,
	addonClient addonclient.Interface,
	addonInformer addoninformerv1alpha1.ManagedClusterAddOnInformer,
	xCMAPIServer string,
	accessToken string,
	recorder events.Recorder,
) factory.Controller {
	controllerName := "addon-controller"
	syncCtx := factory.NewSyncContext(controllerName, recorder)

	c := &addonController{
		clusterLister:             clusterInformer.Lister(),
		addonClient:               addonClient,
		addonLister:               addonInformer.Lister(),
		queue:                     syncCtx.Queue(),
		xCMAPIServer:              xCMAPIServer,
		accessToken:               accessToken,
		addonStatusChangeMap:      make(map[string][]metav1.Condition),
		addonStatusChangeMapMutex: &sync.RWMutex{},
	}

	addonEventHandlerFuncs := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			addon, ok := obj.(*addonv1alpha1.ManagedClusterAddOn)
			if !ok {
				runtime.HandleError(fmt.Errorf("error to get object: %v", obj))
				return
			}
			c.queue.Add(addon.GetNamespace() + "/" + addon.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// only need handle addon status update
			oldAddon, ok := oldObj.(*addonv1alpha1.ManagedClusterAddOn)
			if !ok {
				runtime.HandleError(fmt.Errorf("error to get object: %v", oldObj))
				return
			}
			newAddon, ok := newObj.(*addonv1alpha1.ManagedClusterAddOn)
			if !ok {
				runtime.HandleError(fmt.Errorf("error to get object: %v", newObj))
				return
			}
			if !reflect.DeepEqual(oldAddon.Status.Conditions, newAddon.Status.Conditions) {
				c.updateAddonChangeStatus(newAddon.GetNamespace(), newAddon.GetName(), oldAddon.Status.Conditions, newAddon.Status.Conditions)
				c.queue.Add(newAddon.GetNamespace() + "/" + newAddon.GetName())
			}
		},
		DeleteFunc: func(obj interface{}) {
			switch t := obj.(type) {
			case *addonv1alpha1.ManagedClusterAddOn:
				c.queue.Add(t.GetNamespace() + "/" + t.GetName())
			case cache.DeletedFinalStateUnknown:
				addon, ok := t.Obj.(*addonv1alpha1.ManagedClusterAddOn)
				if !ok {
					runtime.HandleError(fmt.Errorf("error to get object: %v", obj))
					return
				}
				c.queue.Add(addon.GetNamespace() + "/" + addon.GetName())
			default:
				runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
				return
			}
		},
	}

	_, err := addonInformer.Informer().AddEventHandler(addonEventHandlerFuncs)

	if err != nil {
		runtime.HandleError(err)
	}

	return factory.New().
		WithSyncContext(syncCtx).
		// WithInformersQueueKeyFunc(func(obj runtime.Object) string {
		// 	accessor, _ := meta.Accessor(obj)
		// 	return accessor.GetName()
		// }, clusterInformer.Informer()).
		WithBareInformers(addonInformer.Informer()).
		WithSync(c.sync).
		ToController(controllerName, recorder)
}

func (c *addonController) sync(ctx context.Context, syncCtx factory.SyncContext) error {
	queueKey := syncCtx.QueueKey()
	clusterName, addonName, err := cache.SplitMetaNamespaceKey(queueKey)
	if err != nil {
		runtime.HandleError(err)
		return nil
	}

	klog.Infof("Reconciling ManagedClusterAddon %s for Cluster %s", addonName, clusterName)
	addon, err := c.addonLister.ManagedClusterAddOns(clusterName).Get(addonName)
	if errors.IsNotFound(err) {
		// addon not found, could have been deleted, do nothing.
		return nil
	}
	if err != nil {
		return err
	}
	if !addon.DeletionTimestamp.IsZero() {
		// addon is being deleted, do nothing.
		return nil
	}

	cluster, err := c.clusterLister.Get(clusterName)
	if errors.IsNotFound(err) {
		// no cluster, it could be deleted
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to find cluster with name %q: %w", clusterName, err)
	}
	// no work if cluster is being deleted
	if !cluster.DeletionTimestamp.IsZero() {
		return nil
	}

	// check condition firstly
	cond := meta.FindStatusCondition(addon.Status.Conditions, ManagedClusterAddonConditionConnected)
	if cond != nil {
		// update addon status
		c.addonStatusChangeMapMutex.RLock()
		addContidions := c.addonStatusChangeMap[clusterName+"/"+addonName]
		c.addonStatusChangeMapMutex.RUnlock()
		if len(addContidions) > 0 {
			if err := helpers.PostAddonStatus(c.xCMAPIServer, c.accessToken, cluster, addonName, addContidions); err != nil {
				return fmt.Errorf("failed to update status conditions %v for addon %s/%s: %v", addContidions, cluster.GetName(), addonName, err)
			} else {
				addContidions = []metav1.Condition{}
			}
		}

		c.addonStatusChangeMapMutex.Lock()
		// reset the status change
		c.addonStatusChangeMap[clusterName+"/"+addonName] = addContidions
		c.addonStatusChangeMapMutex.Unlock()

		return nil
	}

	// TODO timeout or retry limits
	if err := helpers.CreateAddon(c.xCMAPIServer, c.accessToken, cluster, addon); err != nil {
		return err
	}

	conditionUpdateFn := helpers.UpdateManagedClusterAddonConditionFn(metav1.Condition{
		Type:    ManagedClusterAddonConditionConnected,
		Status:  metav1.ConditionTrue,
		Reason:  "ManagedClusterAddonConnected",
		Message: fmt.Sprintf("Managed cluster addon %s/%s is connected to xCM.", clusterName, addonName),
	})

	_, updated, err := helpers.UpdateManagedClusterAddonStatus(ctx, c.addonClient, clusterName, addonName, conditionUpdateFn)
	if updated {
		syncCtx.Recorder().Eventf("ManagedClusterAddonConnectedConditionUpdated",
			"update managed cluster addon %q/%q connected condition to true",
			clusterName, addonName)
	}

	return err
}

// updateAddonChangeStatus
func (c *addonController) updateAddonChangeStatus(clusterName, addonName string, old, new []metav1.Condition) {
	c.addonStatusChangeMapMutex.RLock()
	addContidions := c.addonStatusChangeMap[clusterName+"/"+addonName]
	c.addonStatusChangeMapMutex.RUnlock()

	for _, cond := range new {
		oldCond := meta.FindStatusCondition(old, cond.Type)
		if oldCond != nil && oldCond.Status == cond.Status {
			continue
		}
		meta.SetStatusCondition(&addContidions, cond)
	}

	c.addonStatusChangeMapMutex.Lock()
	c.addonStatusChangeMap[clusterName+"/"+addonName] = addContidions
	c.addonStatusChangeMapMutex.Unlock()
}
