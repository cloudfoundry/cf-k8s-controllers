package workloads

import (
	"context"
	"errors"

	"code.cloudfoundry.org/cf-k8s-controllers/controllers/apis/workloads/v1alpha1"
	"code.cloudfoundry.org/cf-k8s-controllers/controllers/webhooks"
	"github.com/go-logr/logr"
	admissionv1 "k8s.io/api/admission/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fake -fake-name NameValidator . NameValidator

type NameValidator interface {
	ValidateCreate(ctx context.Context, logger logr.Logger, namespace, newName string) error
	ValidateUpdate(ctx context.Context, logger logr.Logger, namespace, oldName, newName string) error
	ValidateDelete(ctx context.Context, logger logr.Logger, namespace, oldName string) error
}

const (
	AppEntityType = "app"
)

var cfapplog = logf.Log.WithName("cfapp-validate")

//+kubebuilder:webhook:path=/validate-workloads-cloudfoundry-org-v1alpha1-cfapp,mutating=false,failurePolicy=fail,sideEffects=None,groups=workloads.cloudfoundry.org,resources=cfapps,verbs=create;update;delete,versions=v1alpha1,name=vcfapp.workloads.cloudfoundry.org,admissionReviewVersions={v1,v1beta1}

type CFAppValidation struct {
	decoder            *admission.Decoder
	duplicateValidator NameValidator
}

func NewCFAppValidation(duplicateValidator NameValidator) *CFAppValidation {
	return &CFAppValidation{
		duplicateValidator: duplicateValidator,
	}
}

func (v *CFAppValidation) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().Register("/validate-workloads-cloudfoundry-org-v1alpha1-cfapp", &webhook.Admission{Handler: v})

	return nil
}

func (v *CFAppValidation) Handle(ctx context.Context, req admission.Request) admission.Response {
	cfapplog.Info("Validate", "name", req.Name)

	var cfApp, oldCFApp v1alpha1.CFApp
	if req.Operation == admissionv1.Create || req.Operation == admissionv1.Update {
		err := v.decoder.Decode(req, &cfApp)
		if err != nil {
			errMessage := "Error while decoding CFApp object"
			cfapplog.Error(err, errMessage)

			return admission.Denied(errMessage)
		}
	}
	if req.Operation == admissionv1.Update || req.Operation == admissionv1.Delete {
		err := v.decoder.DecodeRaw(req.OldObject, &oldCFApp)
		if err != nil {
			errMessage := "Error while decoding old CFApp object"
			cfapplog.Error(err, errMessage)

			return admission.Denied(errMessage)
		}
	}

	var validatorErr error
	switch req.Operation {
	case admissionv1.Create:
		validatorErr = v.duplicateValidator.ValidateCreate(ctx, cfapplog, cfApp.Namespace, cfApp.Spec.Name)

	case admissionv1.Update:
		validatorErr = v.duplicateValidator.ValidateUpdate(ctx, cfapplog, cfApp.Namespace, oldCFApp.Spec.Name, cfApp.Spec.Name)

	case admissionv1.Delete:
		validatorErr = v.duplicateValidator.ValidateDelete(ctx, cfapplog, oldCFApp.Namespace, oldCFApp.Spec.Name)
	}

	if validatorErr != nil {
		if errors.Is(validatorErr, webhooks.ErrorDuplicateName) {
			return admission.Denied(webhooks.DuplicateAppError.Marshal())
		}

		return admission.Denied(webhooks.UnknownError.Marshal())
	}

	return admission.Allowed("")
}

func (v *CFAppValidation) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
