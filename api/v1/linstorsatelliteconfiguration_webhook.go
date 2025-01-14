/*
Copyright 2022.

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

package v1

import (
	"strconv"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var linstorsatelliteconfigurationlog = logf.Log.WithName("linstorsatelliteconfiguration-resource")

func (r *LinstorSatelliteConfiguration) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-piraeus-io-v1-linstorsatelliteconfiguration,mutating=false,failurePolicy=fail,sideEffects=None,groups=piraeus.io,resources=linstorsatelliteconfigurations,verbs=create;update,versions=v1,name=vlinstorsatelliteconfiguration.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &LinstorSatelliteConfiguration{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LinstorSatelliteConfiguration) ValidateCreate() error {
	linstorsatelliteconfigurationlog.Info("validate create", "name", r.Name)

	errs := r.validate(nil)
	if len(errs) != 0 {
		return apierrors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: "LinstorSatelliteConfiguration"}, r.Name, errs)
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LinstorSatelliteConfiguration) ValidateUpdate(old runtime.Object) error {
	linstorsatelliteconfigurationlog.Info("validate update", "name", r.Name)

	errs := r.validate(old.(*LinstorSatelliteConfiguration))
	if len(errs) != 0 {
		return apierrors.NewInvalid(schema.GroupKind{Group: GroupVersion.Group, Kind: "LinstorSatelliteConfiguration"}, r.Name, errs)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LinstorSatelliteConfiguration) ValidateDelete() error {
	linstorsatelliteconfigurationlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *LinstorSatelliteConfiguration) validate(old *LinstorSatelliteConfiguration) field.ErrorList {
	var oldSPs []LinstorStoragePool
	if old != nil {
		oldSPs = old.Spec.StoragePools
	}

	errs := ValidateStoragePools(r.Spec.StoragePools, oldSPs, field.NewPath("spec", "storagePools"))
	errs = append(errs, ValidateNodeSelector(r.Spec.NodeSelector, field.NewPath("spec", "nodeSelector"))...)
	errs = append(errs, ValidateNodeProperties(r.Spec.Properties, field.NewPath("spec", "properties"))...)

	for i := range r.Spec.Patches {
		errs = append(errs, r.Spec.Patches[i].validate(field.NewPath("spec", "patches", strconv.Itoa(i)))...)
	}

	return errs
}

func ValidateNodeSelector(selector map[string]string, path *field.Path) field.ErrorList {
	var result field.ErrorList

	for k, v := range selector {
		errs := validation.IsQualifiedName(k)
		for _, e := range errs {
			result = append(result, field.Invalid(path, k, e))
		}

		errs = validation.IsValidLabelValue(v)
		for _, e := range errs {
			result = append(result, field.Invalid(path.Child(k), v, e))
		}
	}

	return result
}
