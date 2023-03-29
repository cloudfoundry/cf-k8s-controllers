package env

import (
	"context"
	"encoding/json"
	"fmt"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/shared"
	trinityv1alpha1 "github.tools.sap/neoCoreArchitecture/trinity-service-manager/controllers/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VCAPServicesEnvValueBuilder struct {
	k8sClient client.Client
}

func NewVCAPServicesEnvValueBuilder(k8sClient client.Client) *VCAPServicesEnvValueBuilder {
	return &VCAPServicesEnvValueBuilder{k8sClient: k8sClient}
}

func (b *VCAPServicesEnvValueBuilder) BuildEnvValue(ctx context.Context, cfApp *korifiv1alpha1.CFApp) (map[string]string, error) {
	serviceBindings := &korifiv1alpha1.CFServiceBindingList{}
	err := b.k8sClient.List(ctx, serviceBindings,
		client.InNamespace(cfApp.Namespace),
		client.MatchingFields{shared.IndexServiceBindingAppGUID: cfApp.Name},
	)
	if err != nil {
		return nil, fmt.Errorf("error listing CFServiceBindings: %w", err)
	}

	if len(serviceBindings.Items) == 0 {
		return map[string]string{"VCAP_SERVICES": "{}"}, nil
	}

	vcapServices := map[string][]ServiceDetails{}

	for _, currentServiceBinding := range serviceBindings.Items {
		// If finalizing do not append
		if !currentServiceBinding.DeletionTimestamp.IsZero() {
			continue
		}

		var serviceEnv ServiceDetails
		serviceEnv, err = b.buildSingleServiceEnv(ctx, b.k8sClient, currentServiceBinding)
		if err != nil {
			return nil, err
		}

		if _, ok := vcapServices[serviceEnv.Label]; !ok {
			vcapServices[serviceEnv.Label] = []ServiceDetails{}
		}
		vcapServices[serviceEnv.Label] = append(vcapServices[serviceEnv.Label], serviceEnv)
	}

	jsonVal, err := json.Marshal(vcapServices)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"VCAP_SERVICES": string(jsonVal),
	}, nil
}

func (b *VCAPServicesEnvValueBuilder) getServicePlan(ctx context.Context, servicePlanGuid string) (trinityv1alpha1.CFServicePlan, error) {
	servicePlans := trinityv1alpha1.CFServicePlanList{}
	err := b.k8sClient.List(ctx, &servicePlans, client.MatchingFields{shared.IndexServicePlanGUID: servicePlanGuid})
	if err != nil {
		return trinityv1alpha1.CFServicePlan{}, err
	}

	if len(servicePlans.Items) != 1 {
		return trinityv1alpha1.CFServicePlan{}, fmt.Errorf("found %d service plans for guid %q, expected one", len(servicePlans.Items), servicePlanGuid)
	}

	return servicePlans.Items[0], nil
}

func (b *VCAPServicesEnvValueBuilder) getServiceOffering(ctx context.Context, serviceOfferingGuid string) (trinityv1alpha1.CFServiceOffering, error) {
	serviceOfferings := trinityv1alpha1.CFServiceOfferingList{}
	err := b.k8sClient.List(ctx, &serviceOfferings, client.MatchingFields{shared.IndexServiceOfferingGUID: serviceOfferingGuid})
	if err != nil {
		return trinityv1alpha1.CFServiceOffering{}, err
	}

	if len(serviceOfferings.Items) != 1 {
		return trinityv1alpha1.CFServiceOffering{}, fmt.Errorf("found %d service offerings for guid %q, expected one", len(serviceOfferings.Items), serviceOfferingGuid)
	}

	return serviceOfferings.Items[0], nil
}

func (b *VCAPServicesEnvValueBuilder) buildSingleServiceEnv(ctx context.Context, k8sClient client.Client, serviceBinding korifiv1alpha1.CFServiceBinding) (ServiceDetails, error) {
	if serviceBinding.Status.Binding.Name == "" {
		return ServiceDetails{}, fmt.Errorf("service binding secret name is empty")
	}

	serviceInstance := korifiv1alpha1.CFServiceInstance{}
	err := k8sClient.Get(ctx, types.NamespacedName{Namespace: serviceBinding.Namespace, Name: serviceBinding.Spec.Service.Name}, &serviceInstance)
	if err != nil {
		return ServiceDetails{}, fmt.Errorf("error fetching CFServiceInstance: %w", err)
	}

	secret := corev1.Secret{}
	err = k8sClient.Get(ctx, types.NamespacedName{Namespace: serviceBinding.Namespace, Name: serviceBinding.Status.Binding.Name}, &secret)
	if err != nil {
		return ServiceDetails{}, fmt.Errorf("error fetching CFServiceBinding Secret: %w", err)
	}

	return b.fromServiceBinding(ctx, serviceBinding, serviceInstance, secret)
}

func (b *VCAPServicesEnvValueBuilder) getServiceLabel(ctx context.Context, cfServiceInstance korifiv1alpha1.CFServiceInstance) (string, error) {
	if cfServiceInstance.Spec.Type == korifiv1alpha1.UserProvidedType {
		return korifiv1alpha1.UserProvidedType, nil
	}

	servicePlan, err := b.getServicePlan(ctx, cfServiceInstance.Spec.ServicePlanGUID)
	if err != nil {
		return "", err
	}

	serviceOffering, err := b.getServiceOffering(ctx, servicePlan.Spec.Relationships.ServiceOfferingGUID)
	if err != nil {
		return "", err
	}

	return serviceOffering.Spec.OfferingName, nil
}

func (b *VCAPServicesEnvValueBuilder) fromServiceBinding(
	ctx context.Context,
	serviceBinding korifiv1alpha1.CFServiceBinding,
	serviceInstance korifiv1alpha1.CFServiceInstance,
	serviceBindingSecret corev1.Secret,
) (ServiceDetails, error) {
	var serviceName string
	var bindingName *string

	if serviceBinding.Spec.DisplayName != nil {
		serviceName = *serviceBinding.Spec.DisplayName
		bindingName = serviceBinding.Spec.DisplayName
	} else {
		serviceName = serviceInstance.Spec.DisplayName
		bindingName = nil
	}

	tags := serviceInstance.Spec.Tags
	if tags == nil {
		tags = []string{}
	}
	bindingTagBytes, ok := serviceBindingSecret.Data["tags"]
	if ok {
		var bindingTags []string
		if err := json.Unmarshal(bindingTagBytes, &bindingTags); err != nil {
			return ServiceDetails{}, fmt.Errorf("failed to unmarshal binding secret tags: %w", err)
		}

		tags = append(tags, bindingTags...)
	}

	label, err := b.getServiceLabel(ctx, serviceInstance)
	if err != nil {
		return ServiceDetails{}, err
	}

	credentials, err := mapFromSecret(serviceBindingSecret)
	if err != nil {
		return ServiceDetails{}, err
	}

	return ServiceDetails{
		Label:          label,
		Name:           serviceName,
		Tags:           tags,
		InstanceGUID:   serviceInstance.Name,
		InstanceName:   serviceInstance.Spec.DisplayName,
		BindingGUID:    serviceBinding.Name,
		BindingName:    bindingName,
		Credentials:    credentials,
		SyslogDrainURL: nil,
		VolumeMounts:   []string{},
	}, nil
}

func mapFromSecret(secret corev1.Secret) (map[string]any, error) {
	convertedMap := make(map[string]any)
	for k := range secret.Data {
		var err error
		convertedMap[k], err = parseValue(secret, k)
		if err != nil {
			return nil, err
		}
	}

	return convertedMap, nil
}

type propertyMetadata struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

func parseValue(bindingSecret corev1.Secret, key string) (any, error) {
	valueFormat, err := getValueFormat(bindingSecret, key)
	if err != nil {
		return nil, err
	}

	switch valueFormat {
	case "text":
		return string(bindingSecret.Data[key]), nil
	case "json":
		var value any
		json.Unmarshal(bindingSecret.Data[key], &value)
		return value, nil

	}

	return nil, fmt.Errorf("unsupported value format %q for key %q in secret %s/%s", valueFormat, key, bindingSecret.Namespace, bindingSecret.Name)
}

func getValueFormat(bindingSecret corev1.Secret, key string) (string, error) {
	secretMetadata, ok := bindingSecret.Data[".metadata"]
	if !ok {
		return "text", nil
	}

	var metadata map[string][]propertyMetadata
	if err := json.Unmarshal(secretMetadata, &metadata); err != nil {
		return "", fmt.Errorf("failed to unmarshal metadata from secret %s/%s: %w", bindingSecret.Namespace, bindingSecret.Name, err)
	}

	for _, properties := range metadata {
		if valueFormat := getPropertyFormat(properties, key); valueFormat != "" {
			return valueFormat, nil
		}
	}

	return "text", nil
}

func getPropertyFormat(credentialProperties []propertyMetadata, propertyName string) string {
	for _, property := range credentialProperties {
		if property.Name == propertyName {
			return property.Format
		}
	}
	return ""
}
