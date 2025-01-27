package osbapi

import (
	"context"

	korifiv1alpha1 "code.cloudfoundry.org/korifi/controllers/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//counterfeiter:generate -o fake -fake-name AssetsClient code.cloudfoundry.org/korifi/controllers/controllers/services/osbapi.AssetsClient
type AssetsClient interface {
	GetServiceInstanceAssets(context.Context, *korifiv1alpha1.CFServiceInstance) (ServiceInstanceAssets, error)
	GetServiceBindingAssets(context.Context, *korifiv1alpha1.CFServiceBinding) (ServiceBindingAssets, error)
}

//counterfeiter:generate -o fake -fake-name AssetsClientFactory code.cloudfoundry.org/korifi/controllers/controllers/services/osbapi.AssetsClientFactory
type AssetsClientFactory interface {
	CreateAssetsClient(client.Client, string) AssetsClient
}

type AssetsFactory struct {
}

func NewAssetsFactory() *AssetsFactory {
	return &AssetsFactory{}
}

func (af *AssetsFactory) CreateAssetsClient(k8sClient client.Client, rootNamespace string) AssetsClient {
	return &Assets{
		k8sClient:     k8sClient,
		rootNamespace: rootNamespace,
	}
}
