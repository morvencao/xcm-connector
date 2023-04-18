package helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	operatorhelpers "github.com/openshift/library-go/pkg/operator/v1helpers"

	clusterv1 "open-cluster-management.io/api/cluster/v1"
	// "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusteropenapi "gitlab.cee.redhat.com/service/cluster-managed-service/pkg/api/openapi"
)

func CreateCluster(server, accessToken string, managedCluster *clusterv1.ManagedCluster) error {
	cluster, err := toCluster(managedCluster)
	if err != nil {
		return err
	}

	clusterData, err := json.Marshal(cluster)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/cluster_self_managed/v1/clusters", server)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(clusterData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", fmt.Sprintf("AccessToken %s:%s", *cluster.Id, accessToken))

	// TODO: remove InsecureSkipVerify
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create cluster %s statuscode=%d, status=%s",
			managedCluster.Name, resp.StatusCode, resp.Status)
	}

	return nil
}

func toCluster(managedCluster *clusterv1.ManagedCluster) (*clusteropenapi.Cluster, error) {
	clusterID := findClusterClaims(managedCluster.Status.ClusterClaims, "id.k8s.io")
	if clusterID == "unknown" {
		return nil, fmt.Errorf("failed to get clustrer id for cluster %s", managedCluster.GetName())
	}

	// status := "Unknown"
	// available := meta.FindStatusCondition(managedCluster.Status.Conditions, clusterv1.ManagedClusterConditionAvailable)
	// if available != nil {
	// 	switch available.Status {
	// 	case metav1.ConditionTrue:
	// 		status = "Available"
	// 	case metav1.ConditionFalse:
	// 		status = "Unavailable"
	// 	}
	// }

	return &clusteropenapi.Cluster{
		Kind:           New("Cluster"),
		Id:             &clusterID,
		ControlplaneId: &clusterID,
		// Status:   status,
		// Type:     findClusterClaims(managedCluster.Status.ClusterClaims, "product.open-cluster-management.io"),
		// Version:  managedCluster.Status.Version.Kubernetes,
		// Platform: findClusterClaims(managedCluster.Status.ClusterClaims, "platform.open-cluster-management.io"),
		// Region:   findClusterClaims(managedCluster.Status.ClusterClaims, "region.open-cluster-management.io"),
	}, nil
}

func findClusterClaims(claims []clusterv1.ManagedClusterClaim, name string) string {
	for _, claim := range claims {
		if claim.Name == name {
			return claim.Value
		}
	}

	return "unknown"
}

func PostClusterLabels(server, accessToken string, managedCluster *clusterv1.ManagedCluster, changedLabels map[string]string) error {
	clusterID := findClusterClaims(managedCluster.Status.ClusterClaims, "id.k8s.io")
	if clusterID == "unknown" {
		return fmt.Errorf("failed to get clustrer id for cluster %s", managedCluster.GetName())
	}

	clusterLabelList, err := json.Marshal(toClusterLabelList(changedLabels))
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/cluster_self_managed/v1/clusters/%s/labels", server, clusterID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(clusterLabelList))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", fmt.Sprintf("AccessToken %s:%s", clusterID, accessToken))

	// TODO: remove InsecureSkipVerify
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post labels %s for cluster %s statuscode=%d, status=%s",
			string(clusterLabelList), managedCluster.Name, resp.StatusCode, resp.Status)
	}

	return nil
}

func toClusterLabelList(labels map[string]string) *clusteropenapi.ClusterLabelList {
	clusterLabels := make([]clusteropenapi.ClusterLabel, 0)
	for k, v := range labels {
		clusterLabels = append(clusterLabels, clusteropenapi.ClusterLabel{
			Kind:  New("ClusterLabel"),
			Id:    &k,
			Key:   &k,
			Value: &v,
		})
	}
	return &clusteropenapi.ClusterLabelList{
		Kind:  "ClusterLabelList",
		Items: clusterLabels,
	}
}

func RemoveClusterLabels(server, accessToken string, managedCluster *clusterv1.ManagedCluster, removedLabelKeys []string) error {
	clusterID := findClusterClaims(managedCluster.Status.ClusterClaims, "id.k8s.io")
	if clusterID == "unknown" {
		return fmt.Errorf("failed to get clustrer id for cluster %s", managedCluster.GetName())
	}

	var errList []error
	for _, key := range removedLabelKeys {
		url := fmt.Sprintf("%s/api/cluster_self_managed/v1/clusters/%s/labels/%s", server, clusterID, key)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			errList = append(errList, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("Authorization", fmt.Sprintf("AccessToken %s:%s", clusterID, accessToken))

		// TODO: remove InsecureSkipVerify
		client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		resp, err := client.Do(req)
		if err != nil {
			errList = append(errList, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			errList = append(errList, fmt.Errorf("failed to delete cluster label key %s for cluster %s statuscode=%d, status=%s",
				key, managedCluster.Name, resp.StatusCode, resp.Status))
			continue
		}
	}

	if len(errList) > 0 {
		return operatorhelpers.NewMultiLineAggregate(errList)
	}

	return nil
}

func PostClusterStatus(server, accessToken string, managedCluster *clusterv1.ManagedCluster, addContidions []metav1.Condition) error {
	clusterID := findClusterClaims(managedCluster.Status.ClusterClaims, "id.k8s.io")
	if clusterID == "unknown" {
		return fmt.Errorf("failed to get clustrer id for cluster %s", managedCluster.GetName())
	}

	clusterStatusList, err := json.Marshal(toClusterStatusList(addContidions))
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/cluster_self_managed/v1/clusters/%s/status", server, clusterID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(clusterStatusList))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", fmt.Sprintf("AccessToken %s:%s", clusterID, accessToken))

	// TODO: remove InsecureSkipVerify
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post status %v for cluster %s statuscode=%d, status=%s",
			clusterStatusList, managedCluster.Name, resp.StatusCode, resp.Status)
	}

	return nil
}

func toClusterStatusList(contidions []metav1.Condition) *clusteropenapi.ClusterStatusList {
	clusterStatuses := make([]clusteropenapi.ClusterStatus, len(contidions))
	for i, cond := range contidions {
		clusterStatuses[i] = clusteropenapi.ClusterStatus{
			Type:             clusteropenapi.PtrString(cond.Type),
			State:            clusteropenapi.PtrString(string(cond.Status)),
			Message:          clusteropenapi.PtrString(cond.Message),
			UpdatedTimestamp: presentTime(cond.LastTransitionTime.Time),
		}
	}
	return &clusteropenapi.ClusterStatusList{
		Kind:  "ClusterStatusList",
		Items: clusterStatuses,
	}
}

func presentTime(t time.Time) *time.Time {
	if t.IsZero() {
		return clusteropenapi.PtrTime(time.Time{})
	}

	return clusteropenapi.PtrTime(t.Round(time.Microsecond))
}
