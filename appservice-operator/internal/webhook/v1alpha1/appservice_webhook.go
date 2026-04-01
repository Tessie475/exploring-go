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

package v1alpha1

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	appsv1alpha1 "example.com/appservice-operator/api/v1alpha1"
)

// nolint:unused
// log is for logging in this package.
var appservicelog = logf.Log.WithName("appservice-resource")

// SetupAppServiceWebhookWithManager registers the webhook for AppService in the manager.
func SetupAppServiceWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &appsv1alpha1.AppService{}).
		WithValidator(&AppServiceCustomValidator{}).
		WithDefaulter(&AppServiceCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-apps-example-com-v1alpha1-appservice,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.example.com,resources=appservices,verbs=create;update,versions=v1alpha1,name=mappservice-v1alpha1.kb.io,admissionReviewVersions=v1

// AppServiceCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind AppService when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type AppServiceCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind AppService.
// This is the MUTATING webhook — it modifies the resource before it's saved.
func (d *AppServiceCustomDefaulter) Default(_ context.Context, obj *appsv1alpha1.AppService) error {
	appservicelog.Info("Defaulting for AppService", "name", obj.GetName())

	// If no replicas specified (zero value), default to 2
	if obj.Spec.Replicas == 0 {
		obj.Spec.Replicas = 2
		appservicelog.Info("Defaulted replicas to 2")
	}

	// Add a managed-by label if not already present
	if obj.Labels == nil {
		obj.Labels = make(map[string]string)
	}
	if _, exists := obj.Labels["managed-by"]; !exists {
		obj.Labels["managed-by"] = "appservice-operator"
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: If you want to customise the 'path', use the flags '--defaulting-path' or '--validation-path'.
// +kubebuilder:webhook:path=/validate-apps-example-com-v1alpha1-appservice,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.example.com,resources=appservices,verbs=create;update,versions=v1alpha1,name=vappservice-v1alpha1.kb.io,admissionReviewVersions=v1

// AppServiceCustomValidator struct is responsible for validating the AppService resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type AppServiceCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type AppService.
// This is the VALIDATING webhook — it accepts or rejects the resource.
func (v *AppServiceCustomValidator) ValidateCreate(_ context.Context, obj *appsv1alpha1.AppService) (admission.Warnings, error) {
	appservicelog.Info("Validation for AppService upon creation", "name", obj.GetName())

	return validateAppService(obj)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type AppService.
func (v *AppServiceCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj *appsv1alpha1.AppService) (admission.Warnings, error) {
	appservicelog.Info("Validation for AppService upon update", "name", newObj.GetName())

	// Same validation rules apply on update
	return validateAppService(newObj)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type AppService.
func (v *AppServiceCustomValidator) ValidateDelete(_ context.Context, obj *appsv1alpha1.AppService) (admission.Warnings, error) {
	appservicelog.Info("Validation for AppService upon deletion", "name", obj.GetName())

	// Allow all deletions
	return nil, nil
}

// validateAppService contains the shared validation rules used by both create and update.
func validateAppService(obj *appsv1alpha1.AppService) (admission.Warnings, error) {
	var warnings admission.Warnings

	// Image is required
	if obj.Spec.Image == "" {
		return nil, fmt.Errorf("spec.image is required and cannot be empty")
	}

	// Replicas must be between 1 and 10
	if obj.Spec.Replicas < 1 || obj.Spec.Replicas > 10 {
		return nil, fmt.Errorf("spec.replicas must be between 1 and 10, got %d", obj.Spec.Replicas)
	}

	// Port must be valid
	if obj.Spec.Port < 1 || obj.Spec.Port > 65535 {
		return nil, fmt.Errorf("spec.port must be between 1 and 65535, got %d", obj.Spec.Port)
	}

	// Warn (but allow) if using the 'latest' tag — not recommended in production
	if len(obj.Spec.Image) > 7 && obj.Spec.Image[len(obj.Spec.Image)-7:] == ":latest" {
		warnings = append(warnings, "using ':latest' tag is not recommended for production")
	}

	return warnings, nil
}
