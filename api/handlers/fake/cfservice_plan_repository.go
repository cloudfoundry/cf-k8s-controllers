// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFServicePlanRepository struct {
	ListPlansStub        func(context.Context, authorization.Info, repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error)
	listPlansMutex       sync.RWMutex
	listPlansArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListServicePlanMessage
	}
	listPlansReturns struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}
	listPlansReturnsOnCall map[int]struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFServicePlanRepository) ListPlans(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error) {
	fake.listPlansMutex.Lock()
	ret, specificReturn := fake.listPlansReturnsOnCall[len(fake.listPlansArgsForCall)]
	fake.listPlansArgsForCall = append(fake.listPlansArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListServicePlanMessage
	}{arg1, arg2, arg3})
	stub := fake.ListPlansStub
	fakeReturns := fake.listPlansReturns
	fake.recordInvocation("ListPlans", []interface{}{arg1, arg2, arg3})
	fake.listPlansMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFServicePlanRepository) ListPlansCallCount() int {
	fake.listPlansMutex.RLock()
	defer fake.listPlansMutex.RUnlock()
	return len(fake.listPlansArgsForCall)
}

func (fake *CFServicePlanRepository) ListPlansCalls(stub func(context.Context, authorization.Info, repositories.ListServicePlanMessage) ([]repositories.ServicePlanRecord, error)) {
	fake.listPlansMutex.Lock()
	defer fake.listPlansMutex.Unlock()
	fake.ListPlansStub = stub
}

func (fake *CFServicePlanRepository) ListPlansArgsForCall(i int) (context.Context, authorization.Info, repositories.ListServicePlanMessage) {
	fake.listPlansMutex.RLock()
	defer fake.listPlansMutex.RUnlock()
	argsForCall := fake.listPlansArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFServicePlanRepository) ListPlansReturns(result1 []repositories.ServicePlanRecord, result2 error) {
	fake.listPlansMutex.Lock()
	defer fake.listPlansMutex.Unlock()
	fake.ListPlansStub = nil
	fake.listPlansReturns = struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) ListPlansReturnsOnCall(i int, result1 []repositories.ServicePlanRecord, result2 error) {
	fake.listPlansMutex.Lock()
	defer fake.listPlansMutex.Unlock()
	fake.ListPlansStub = nil
	if fake.listPlansReturnsOnCall == nil {
		fake.listPlansReturnsOnCall = make(map[int]struct {
			result1 []repositories.ServicePlanRecord
			result2 error
		})
	}
	fake.listPlansReturnsOnCall[i] = struct {
		result1 []repositories.ServicePlanRecord
		result2 error
	}{result1, result2}
}

func (fake *CFServicePlanRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listPlansMutex.RLock()
	defer fake.listPlansMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFServicePlanRepository) recordInvocation(key string, args []interface{}) {
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

var _ handlers.CFServicePlanRepository = new(CFServicePlanRepository)
