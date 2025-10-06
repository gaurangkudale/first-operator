/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package demovolume

import (
	"context"
	"fmt"
	"time"

	// v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	demovolumev1alpha1 "github.com/gaurangkudale/first-operator/api/demovolume/v1alpha1"
	"github.com/go-logr/logr"
)

// DemovolumeReconciler reconciles a Demovolume object
type DemovolumeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=demovolume.gaurangkudale.github.io,resources=demovolumes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=demovolume.gaurangkudale.github.io,resources=demovolumes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=demovolume.gaurangkudale.github.io,resources=demovolumes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Demovolume object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *DemovolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// TODO(user): your logic here

	log.Info("Enter Reconcile r: ", "req: ", req)
	volume := &demovolumev1alpha1.Demovolume{}
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, volume); err != nil {
		if errors.IsNotFound(err) {
			// The resource was deleted after the reconcile request—return and don't requeue
			return ctrl.Result{}, nil
		}
		// Requeue if there was another error
		return ctrl.Result{}, err
	}
	log.Info("Enter Reconcile: ", "spec:", volume.Spec, "status:", volume.Status, "Kind: ", volume.Kind)
	// Reconcile PVC
	if err := r.reconcilePVC(ctx, volume, log); err != nil {
		log.Error(err, "Failed to reconcile PVC")
		return ctrl.Result{}, err
	}
	if volume.Spec.Name != volume.Status.Name && volume.Spec.Size != volume.Status.Size {
		volume.Status.Name = volume.Spec.Name
		volume.Status.Size = volume.Spec.Size
		if err := r.Status().Update(ctx, volume); err != nil {
			return ctrl.Result{}, err
		}
	}
	const RequeueInterval = 60 * time.Second
	return ctrl.Result{RequeueAfter: RequeueInterval}, nil
}

func (r *DemovolumeReconciler) reconcilePVC(ctx context.Context, volume *demovolumev1alpha1.Demovolume, l logr.Logger) error {
	pvc := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      volume.Spec.Name,
		Namespace: volume.Namespace,
	}, pvc)

	if err == nil {
		// PVC already exists
		l.Info("PVC already exists", "pvc", pvc.Name)
		return nil
	}

	if !errors.IsNotFound(err) {
		// Other error fetching PVC
		return err
	}

	// PVC not found — let's create it
	storageQty := resource.MustParse(fmt.Sprintf("%dGi", volume.Spec.Size))

	pvc = &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      volume.Spec.Name,
			Namespace: volume.Namespace,
			Labels: map[string]string{
				"app": "demovolume",
				"cr":  volume.Name,
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: storageQty,
				},
			},
		},
	}
	// Set owner reference so PVC is deleted when Demovolume is deleted
	if err := ctrl.SetControllerReference(volume, pvc, r.Scheme); err != nil {
		return err
	}

	// Create the PVC
	if err := r.Create(ctx, pvc); err != nil {
		return err
	}

	l.Info("Created new PVC", "pvc", pvc.Name)
	const RequeueInterval = 60 * time.Second
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemovolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demovolumev1alpha1.Demovolume{}).
		Named("demovolume-demovolume").
		Complete(r)
}
