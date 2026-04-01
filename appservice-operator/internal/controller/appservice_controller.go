/*
Copyright 2026.

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

package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "example.com/appservice-operator/api/v1alpha1"
)

// AppServiceReconciler reconciles a AppService object
type AppServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// RBAC markers tell Kubebuilder what permissions the operator needs.
// We need access to AppServices (our CRD) plus Deployments and Services (what we create).

// +kubebuilder:rbac:groups=apps.example.com,resources=appservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=appservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.example.com,resources=appservices/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile runs every time an AppService is created, updated, or deleted.
// Its job: make the actual cluster state match the desired state in the spec.
func (r *AppServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// ──────────────────────────────────────────────
	// Step 1: Fetch the AppService resource
	// ──────────────────────────────────────────────
	// req.NamespacedName contains the name and namespace of the resource that triggered this reconcile.
	var appService appsv1alpha1.AppService
	if err := r.Get(ctx, req.NamespacedName, &appService); err != nil {
		if apierrors.IsNotFound(err) {
			// The AppService was deleted — nothing to do.
			// Any owned Deployments/Services get auto-deleted via owner references.
			log.Info("AppService not found, must have been deleted")
			return ctrl.Result{}, nil
		}
		// Some other error (network issue, permission problem, etc.)
		return ctrl.Result{}, fmt.Errorf("fetching AppService: %w", err)
	}

	log.Info("Reconciling AppService", "name", appService.Name)

	// ──────────────────────────────────────────────
	// Step 2: Create or Update the Deployment
	// ──────────────────────────────────────────────
	// We build the Deployment object and use CreateOrUpdate to apply it.
	// CreateOrUpdate checks if it exists — if yes, it updates; if no, it creates.

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appService.Name,
			Namespace: appService.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		// This function is called inside CreateOrUpdate to set the desired state.
		// It runs both on create (empty object) and update (existing object).

		replicas := appService.Spec.Replicas
		labels := map[string]string{
			"app":        appService.Name,
			"managed-by": "appservice-operator",
		}

		// Set the Deployment spec to match our AppService spec
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  appService.Name,
							Image: appService.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: appService.Spec.Port,
								},
							},
						},
					},
				},
			},
		}

		// SetControllerReference marks the Deployment as "owned by" this AppService.
		// When the AppService is deleted, Kubernetes auto-deletes the Deployment too.
		return ctrl.SetControllerReference(&appService, deployment, r.Scheme)
	})

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("creating/updating Deployment: %w", err)
	}
	log.Info("Deployment reconciled", "name", deployment.Name, "result", result)

	// ──────────────────────────────────────────────
	// Step 3: Create or Update the Service
	// ──────────────────────────────────────────────
	// The Service routes network traffic to the pods created by the Deployment.

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appService.Name,
			Namespace: appService.Namespace,
		},
	}

	result, err = controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		service.Spec = corev1.ServiceSpec{
			Selector: map[string]string{
				"app": appService.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       appService.Spec.Port,
					TargetPort: intstr.FromInt32(appService.Spec.Port),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		}

		return ctrl.SetControllerReference(&appService, service, r.Scheme)
	})

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("creating/updating Service: %w", err)
	}
	log.Info("Service reconciled", "name", service.Name, "result", result)

	// ──────────────────────────────────────────────
	// Step 4: Update the AppService status
	// ──────────────────────────────────────────────
	// Read the Deployment's actual state and write it back to our status.
	// This lets users see how many replicas are actually running.

	var currentDeployment appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKeyFromObject(deployment), &currentDeployment); err != nil {
		return ctrl.Result{}, fmt.Errorf("fetching Deployment status: %w", err)
	}

	appService.Status.AvailableReplicas = currentDeployment.Status.AvailableReplicas

	// Use Status().Update() — NOT regular Update() — because status is a subresource.
	// This only updates the status fields, not the spec.
	if err := r.Status().Update(ctx, &appService); err != nil {
		return ctrl.Result{}, fmt.Errorf("updating AppService status: %w", err)
	}

	log.Info("Status updated", "availableReplicas", appService.Status.AvailableReplicas)

	return ctrl.Result{}, nil
}

// SetupWithManager registers the controller with the manager and tells it
// to watch AppService resources AND owned Deployments and Services.
func (r *AppServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.AppService{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named("appservice").
		Complete(r)
}
