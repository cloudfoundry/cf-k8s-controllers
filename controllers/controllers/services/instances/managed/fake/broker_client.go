// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/controllers/controllers/services/osbapi"
)

type BrokerClient struct {
	BindStub        func(context.Context, osbapi.BindPayload) (osbapi.BindResponse, error)
	bindMutex       sync.RWMutex
	bindArgsForCall []struct {
		arg1 context.Context
		arg2 osbapi.BindPayload
	}
	bindReturns struct {
		result1 osbapi.BindResponse
		result2 error
	}
	bindReturnsOnCall map[int]struct {
		result1 osbapi.BindResponse
		result2 error
	}
	DeprovisionStub        func(context.Context, osbapi.InstanceDeprovisionPayload) (osbapi.ServiceInstanceOperationResponse, error)
	deprovisionMutex       sync.RWMutex
	deprovisionArgsForCall []struct {
		arg1 context.Context
		arg2 osbapi.InstanceDeprovisionPayload
	}
	deprovisionReturns struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}
	deprovisionReturnsOnCall map[int]struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}
	GetCatalogStub        func(context.Context) (osbapi.Catalog, error)
	getCatalogMutex       sync.RWMutex
	getCatalogArgsForCall []struct {
		arg1 context.Context
	}
	getCatalogReturns struct {
		result1 osbapi.Catalog
		result2 error
	}
	getCatalogReturnsOnCall map[int]struct {
		result1 osbapi.Catalog
		result2 error
	}
	GetServiceInstanceLastOperationStub        func(context.Context, osbapi.GetLastOperationRequest) (osbapi.LastOperationResponse, error)
	getServiceInstanceLastOperationMutex       sync.RWMutex
	getServiceInstanceLastOperationArgsForCall []struct {
		arg1 context.Context
		arg2 osbapi.GetLastOperationRequest
	}
	getServiceInstanceLastOperationReturns struct {
		result1 osbapi.LastOperationResponse
		result2 error
	}
	getServiceInstanceLastOperationReturnsOnCall map[int]struct {
		result1 osbapi.LastOperationResponse
		result2 error
	}
	ProvisionStub        func(context.Context, osbapi.InstanceProvisionPayload) (osbapi.ServiceInstanceOperationResponse, error)
	provisionMutex       sync.RWMutex
	provisionArgsForCall []struct {
		arg1 context.Context
		arg2 osbapi.InstanceProvisionPayload
	}
	provisionReturns struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}
	provisionReturnsOnCall map[int]struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *BrokerClient) Bind(arg1 context.Context, arg2 osbapi.BindPayload) (osbapi.BindResponse, error) {
	fake.bindMutex.Lock()
	ret, specificReturn := fake.bindReturnsOnCall[len(fake.bindArgsForCall)]
	fake.bindArgsForCall = append(fake.bindArgsForCall, struct {
		arg1 context.Context
		arg2 osbapi.BindPayload
	}{arg1, arg2})
	stub := fake.BindStub
	fakeReturns := fake.bindReturns
	fake.recordInvocation("Bind", []interface{}{arg1, arg2})
	fake.bindMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BrokerClient) BindCallCount() int {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	return len(fake.bindArgsForCall)
}

func (fake *BrokerClient) BindCalls(stub func(context.Context, osbapi.BindPayload) (osbapi.BindResponse, error)) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = stub
}

func (fake *BrokerClient) BindArgsForCall(i int) (context.Context, osbapi.BindPayload) {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	argsForCall := fake.bindArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *BrokerClient) BindReturns(result1 osbapi.BindResponse, result2 error) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = nil
	fake.bindReturns = struct {
		result1 osbapi.BindResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) BindReturnsOnCall(i int, result1 osbapi.BindResponse, result2 error) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = nil
	if fake.bindReturnsOnCall == nil {
		fake.bindReturnsOnCall = make(map[int]struct {
			result1 osbapi.BindResponse
			result2 error
		})
	}
	fake.bindReturnsOnCall[i] = struct {
		result1 osbapi.BindResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) Deprovision(arg1 context.Context, arg2 osbapi.InstanceDeprovisionPayload) (osbapi.ServiceInstanceOperationResponse, error) {
	fake.deprovisionMutex.Lock()
	ret, specificReturn := fake.deprovisionReturnsOnCall[len(fake.deprovisionArgsForCall)]
	fake.deprovisionArgsForCall = append(fake.deprovisionArgsForCall, struct {
		arg1 context.Context
		arg2 osbapi.InstanceDeprovisionPayload
	}{arg1, arg2})
	stub := fake.DeprovisionStub
	fakeReturns := fake.deprovisionReturns
	fake.recordInvocation("Deprovision", []interface{}{arg1, arg2})
	fake.deprovisionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BrokerClient) DeprovisionCallCount() int {
	fake.deprovisionMutex.RLock()
	defer fake.deprovisionMutex.RUnlock()
	return len(fake.deprovisionArgsForCall)
}

func (fake *BrokerClient) DeprovisionCalls(stub func(context.Context, osbapi.InstanceDeprovisionPayload) (osbapi.ServiceInstanceOperationResponse, error)) {
	fake.deprovisionMutex.Lock()
	defer fake.deprovisionMutex.Unlock()
	fake.DeprovisionStub = stub
}

func (fake *BrokerClient) DeprovisionArgsForCall(i int) (context.Context, osbapi.InstanceDeprovisionPayload) {
	fake.deprovisionMutex.RLock()
	defer fake.deprovisionMutex.RUnlock()
	argsForCall := fake.deprovisionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *BrokerClient) DeprovisionReturns(result1 osbapi.ServiceInstanceOperationResponse, result2 error) {
	fake.deprovisionMutex.Lock()
	defer fake.deprovisionMutex.Unlock()
	fake.DeprovisionStub = nil
	fake.deprovisionReturns = struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) DeprovisionReturnsOnCall(i int, result1 osbapi.ServiceInstanceOperationResponse, result2 error) {
	fake.deprovisionMutex.Lock()
	defer fake.deprovisionMutex.Unlock()
	fake.DeprovisionStub = nil
	if fake.deprovisionReturnsOnCall == nil {
		fake.deprovisionReturnsOnCall = make(map[int]struct {
			result1 osbapi.ServiceInstanceOperationResponse
			result2 error
		})
	}
	fake.deprovisionReturnsOnCall[i] = struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) GetCatalog(arg1 context.Context) (osbapi.Catalog, error) {
	fake.getCatalogMutex.Lock()
	ret, specificReturn := fake.getCatalogReturnsOnCall[len(fake.getCatalogArgsForCall)]
	fake.getCatalogArgsForCall = append(fake.getCatalogArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.GetCatalogStub
	fakeReturns := fake.getCatalogReturns
	fake.recordInvocation("GetCatalog", []interface{}{arg1})
	fake.getCatalogMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BrokerClient) GetCatalogCallCount() int {
	fake.getCatalogMutex.RLock()
	defer fake.getCatalogMutex.RUnlock()
	return len(fake.getCatalogArgsForCall)
}

func (fake *BrokerClient) GetCatalogCalls(stub func(context.Context) (osbapi.Catalog, error)) {
	fake.getCatalogMutex.Lock()
	defer fake.getCatalogMutex.Unlock()
	fake.GetCatalogStub = stub
}

func (fake *BrokerClient) GetCatalogArgsForCall(i int) context.Context {
	fake.getCatalogMutex.RLock()
	defer fake.getCatalogMutex.RUnlock()
	argsForCall := fake.getCatalogArgsForCall[i]
	return argsForCall.arg1
}

func (fake *BrokerClient) GetCatalogReturns(result1 osbapi.Catalog, result2 error) {
	fake.getCatalogMutex.Lock()
	defer fake.getCatalogMutex.Unlock()
	fake.GetCatalogStub = nil
	fake.getCatalogReturns = struct {
		result1 osbapi.Catalog
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) GetCatalogReturnsOnCall(i int, result1 osbapi.Catalog, result2 error) {
	fake.getCatalogMutex.Lock()
	defer fake.getCatalogMutex.Unlock()
	fake.GetCatalogStub = nil
	if fake.getCatalogReturnsOnCall == nil {
		fake.getCatalogReturnsOnCall = make(map[int]struct {
			result1 osbapi.Catalog
			result2 error
		})
	}
	fake.getCatalogReturnsOnCall[i] = struct {
		result1 osbapi.Catalog
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) GetServiceInstanceLastOperation(arg1 context.Context, arg2 osbapi.GetLastOperationRequest) (osbapi.LastOperationResponse, error) {
	fake.getServiceInstanceLastOperationMutex.Lock()
	ret, specificReturn := fake.getServiceInstanceLastOperationReturnsOnCall[len(fake.getServiceInstanceLastOperationArgsForCall)]
	fake.getServiceInstanceLastOperationArgsForCall = append(fake.getServiceInstanceLastOperationArgsForCall, struct {
		arg1 context.Context
		arg2 osbapi.GetLastOperationRequest
	}{arg1, arg2})
	stub := fake.GetServiceInstanceLastOperationStub
	fakeReturns := fake.getServiceInstanceLastOperationReturns
	fake.recordInvocation("GetServiceInstanceLastOperation", []interface{}{arg1, arg2})
	fake.getServiceInstanceLastOperationMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BrokerClient) GetServiceInstanceLastOperationCallCount() int {
	fake.getServiceInstanceLastOperationMutex.RLock()
	defer fake.getServiceInstanceLastOperationMutex.RUnlock()
	return len(fake.getServiceInstanceLastOperationArgsForCall)
}

func (fake *BrokerClient) GetServiceInstanceLastOperationCalls(stub func(context.Context, osbapi.GetLastOperationRequest) (osbapi.LastOperationResponse, error)) {
	fake.getServiceInstanceLastOperationMutex.Lock()
	defer fake.getServiceInstanceLastOperationMutex.Unlock()
	fake.GetServiceInstanceLastOperationStub = stub
}

func (fake *BrokerClient) GetServiceInstanceLastOperationArgsForCall(i int) (context.Context, osbapi.GetLastOperationRequest) {
	fake.getServiceInstanceLastOperationMutex.RLock()
	defer fake.getServiceInstanceLastOperationMutex.RUnlock()
	argsForCall := fake.getServiceInstanceLastOperationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *BrokerClient) GetServiceInstanceLastOperationReturns(result1 osbapi.LastOperationResponse, result2 error) {
	fake.getServiceInstanceLastOperationMutex.Lock()
	defer fake.getServiceInstanceLastOperationMutex.Unlock()
	fake.GetServiceInstanceLastOperationStub = nil
	fake.getServiceInstanceLastOperationReturns = struct {
		result1 osbapi.LastOperationResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) GetServiceInstanceLastOperationReturnsOnCall(i int, result1 osbapi.LastOperationResponse, result2 error) {
	fake.getServiceInstanceLastOperationMutex.Lock()
	defer fake.getServiceInstanceLastOperationMutex.Unlock()
	fake.GetServiceInstanceLastOperationStub = nil
	if fake.getServiceInstanceLastOperationReturnsOnCall == nil {
		fake.getServiceInstanceLastOperationReturnsOnCall = make(map[int]struct {
			result1 osbapi.LastOperationResponse
			result2 error
		})
	}
	fake.getServiceInstanceLastOperationReturnsOnCall[i] = struct {
		result1 osbapi.LastOperationResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) Provision(arg1 context.Context, arg2 osbapi.InstanceProvisionPayload) (osbapi.ServiceInstanceOperationResponse, error) {
	fake.provisionMutex.Lock()
	ret, specificReturn := fake.provisionReturnsOnCall[len(fake.provisionArgsForCall)]
	fake.provisionArgsForCall = append(fake.provisionArgsForCall, struct {
		arg1 context.Context
		arg2 osbapi.InstanceProvisionPayload
	}{arg1, arg2})
	stub := fake.ProvisionStub
	fakeReturns := fake.provisionReturns
	fake.recordInvocation("Provision", []interface{}{arg1, arg2})
	fake.provisionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *BrokerClient) ProvisionCallCount() int {
	fake.provisionMutex.RLock()
	defer fake.provisionMutex.RUnlock()
	return len(fake.provisionArgsForCall)
}

func (fake *BrokerClient) ProvisionCalls(stub func(context.Context, osbapi.InstanceProvisionPayload) (osbapi.ServiceInstanceOperationResponse, error)) {
	fake.provisionMutex.Lock()
	defer fake.provisionMutex.Unlock()
	fake.ProvisionStub = stub
}

func (fake *BrokerClient) ProvisionArgsForCall(i int) (context.Context, osbapi.InstanceProvisionPayload) {
	fake.provisionMutex.RLock()
	defer fake.provisionMutex.RUnlock()
	argsForCall := fake.provisionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *BrokerClient) ProvisionReturns(result1 osbapi.ServiceInstanceOperationResponse, result2 error) {
	fake.provisionMutex.Lock()
	defer fake.provisionMutex.Unlock()
	fake.ProvisionStub = nil
	fake.provisionReturns = struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) ProvisionReturnsOnCall(i int, result1 osbapi.ServiceInstanceOperationResponse, result2 error) {
	fake.provisionMutex.Lock()
	defer fake.provisionMutex.Unlock()
	fake.ProvisionStub = nil
	if fake.provisionReturnsOnCall == nil {
		fake.provisionReturnsOnCall = make(map[int]struct {
			result1 osbapi.ServiceInstanceOperationResponse
			result2 error
		})
	}
	fake.provisionReturnsOnCall[i] = struct {
		result1 osbapi.ServiceInstanceOperationResponse
		result2 error
	}{result1, result2}
}

func (fake *BrokerClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	fake.deprovisionMutex.RLock()
	defer fake.deprovisionMutex.RUnlock()
	fake.getCatalogMutex.RLock()
	defer fake.getCatalogMutex.RUnlock()
	fake.getServiceInstanceLastOperationMutex.RLock()
	defer fake.getServiceInstanceLastOperationMutex.RUnlock()
	fake.provisionMutex.RLock()
	defer fake.provisionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *BrokerClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ osbapi.BrokerClient = new(BrokerClient)
