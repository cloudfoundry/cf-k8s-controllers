package repositories

import (
	"context"
	"fmt"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"code.cloudfoundry.org/korifi/controllers/controllers/services/osbapi"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ParamsResult struct {
	Parameters map[string]any
}

type BrokerParamsClient struct {
	clientFactory osbapi.BrokerClientFactory
	k8sClient     client.Client
	rootNamespace string
}

func NewBrokerParamsClient(clientFactory osbapi.BrokerClientFactory, k8sClient client.Client, rootNamespace string) *BrokerParamsClient {
	return &BrokerParamsClient{
		clientFactory: clientFactory,
		k8sClient:     k8sClient,
		rootNamespace: rootNamespace,
	}
}

func (c *BrokerParamsClient) GetServiceBindingParameters(ctx context.Context, serviceBinding *korifiv1alpha1.CFServiceBinding) (ParamsResult, error) {
	assetsClient := osbapi.NewAssets(c.k8sClient, c.rootNamespace)
	sbAssets, err := assetsClient.GetServiceBindingAssets(ctx, serviceBinding)
	if err != nil {
		return ParamsResult{}, fmt.Errorf("faild to get service binding assets: %w", err)
	}

	osbapiClient, err := c.clientFactory.CreateClient(ctx, sbAssets.ServiceBroker)
	if err != nil {
		return ParamsResult{}, fmt.Errorf("faild to create osbapi client: %w", err)
	}

	payload := osbapi.BindPayload{
		BindingID:  serviceBinding.Name,
		InstanceID: serviceBinding.Spec.Service.Name,
	}

	binding, err := osbapiClient.GetServiceBinding(ctx, payload)
	if err != nil {
		return ParamsResult{}, fmt.Errorf("faild to create osbapi client: %w", err)
	}

	return ParamsResult{
		Parameters: binding.Parameters,
	}, nil
}
