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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	demovolumev1alpha1 "github.com/gaurangkudale/first-operator/api/demovolume/v1alpha1"
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

	if volume.Spec.Name != volume.Status.Name && volume.Spec.Size != volume.Status.Size {
		volume.Status.Name = volume.Spec.Name
		volume.Status.Size = volume.Spec.Size
		if err := r.Status().Update(ctx, volume); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemovolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demovolumev1alpha1.Demovolume{}).
		Named("demovolume-demovolume").
		Complete(r)
}
