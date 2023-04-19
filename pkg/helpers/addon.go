package helpers

import (
	"context"
	"encoding/json"
	"fmt"

	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	addonclient "open-cluster-management.io/api/client/addon/clientset/versioned"

	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
)

type UpdateManagedClusterAddonStatusFunc func(status *addonv1alpha1.ManagedClusterAddOnStatus) error

func UpdateManagedClusterAddonStatus(
	ctx context.Context,
	client addonclient.Interface,
	clusterName string,
	addonName string,
	updateFuncs ...UpdateManagedClusterAddonStatusFunc) (*addonv1alpha1.ManagedClusterAddOnStatus, bool, error) {
	updated := false
	var updatedManagedClusterAddonStatus *addonv1alpha1.ManagedClusterAddOnStatus

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		addon, err := client.AddonV1alpha1().ManagedClusterAddOns(clusterName).Get(ctx, addonName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		oldStatus := &addon.Status

		newStatus := oldStatus.DeepCopy()
		for _, update := range updateFuncs {
			if err := update(newStatus); err != nil {
				return err
			}
		}
		if equality.Semantic.DeepEqual(oldStatus, newStatus) {
			// We return the newStatus which is a deep copy of oldStatus but with all update funcs applied.
			updatedManagedClusterAddonStatus = newStatus
			return nil
		}

		oldData, err := json.Marshal(addonv1alpha1.ManagedClusterAddOn{
			Status: *oldStatus,
		})

		if err != nil {
			return fmt.Errorf("failed to Marshal old data for addon status %s/%s: %w", clusterName, addonName, err)
		}

		newData, err := json.Marshal(addonv1alpha1.ManagedClusterAddOn{
			ObjectMeta: metav1.ObjectMeta{
				UID:             addon.UID,
				ResourceVersion: addon.ResourceVersion,
			}, // to ensure they appear in the patch as preconditions
			Status: *newStatus,
		})
		if err != nil {
			return fmt.Errorf("failed to Marshal new data for addon status %s/%s: %w", clusterName, addonName, err)
		}

		patchBytes, err := jsonpatch.CreateMergePatch(oldData, newData)
		if err != nil {
			return fmt.Errorf("failed to create patch for addon %s/%s: %w", clusterName, addonName, err)
		}

		updatedManagedClusterAddon, err := client.AddonV1alpha1().ManagedClusterAddOns(clusterName).Patch(ctx, addonName, types.MergePatchType, patchBytes, metav1.PatchOptions{}, "status")

		updatedManagedClusterAddonStatus = &updatedManagedClusterAddon.Status
		updated = err == nil
		return err
	})

	return updatedManagedClusterAddonStatus, updated, err
}

func UpdateManagedClusterAddonConditionFn(cond metav1.Condition) UpdateManagedClusterAddonStatusFunc {
	return func(oldStatus *addonv1alpha1.ManagedClusterAddOnStatus) error {
		meta.SetStatusCondition(&oldStatus.Conditions, cond)
		return nil
	}
}
