// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/apis"
	"code.cloudfoundry.org/korifi/api/authorization"
)

type IdentityProvider struct {
	GetIdentityStub        func(context.Context, authorization.Info) (authorization.Identity, error)
	getIdentityMutex       sync.RWMutex
	getIdentityArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
	}
	getIdentityReturns struct {
		result1 authorization.Identity
		result2 error
	}
	getIdentityReturnsOnCall map[int]struct {
		result1 authorization.Identity
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *IdentityProvider) GetIdentity(arg1 context.Context, arg2 authorization.Info) (authorization.Identity, error) {
	fake.getIdentityMutex.Lock()
	ret, specificReturn := fake.getIdentityReturnsOnCall[len(fake.getIdentityArgsForCall)]
	fake.getIdentityArgsForCall = append(fake.getIdentityArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
	}{arg1, arg2})
	stub := fake.GetIdentityStub
	fakeReturns := fake.getIdentityReturns
	fake.recordInvocation("GetIdentity", []interface{}{arg1, arg2})
	fake.getIdentityMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *IdentityProvider) GetIdentityCallCount() int {
	fake.getIdentityMutex.RLock()
	defer fake.getIdentityMutex.RUnlock()
	return len(fake.getIdentityArgsForCall)
}

func (fake *IdentityProvider) GetIdentityCalls(stub func(context.Context, authorization.Info) (authorization.Identity, error)) {
	fake.getIdentityMutex.Lock()
	defer fake.getIdentityMutex.Unlock()
	fake.GetIdentityStub = stub
}

func (fake *IdentityProvider) GetIdentityArgsForCall(i int) (context.Context, authorization.Info) {
	fake.getIdentityMutex.RLock()
	defer fake.getIdentityMutex.RUnlock()
	argsForCall := fake.getIdentityArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *IdentityProvider) GetIdentityReturns(result1 authorization.Identity, result2 error) {
	fake.getIdentityMutex.Lock()
	defer fake.getIdentityMutex.Unlock()
	fake.GetIdentityStub = nil
	fake.getIdentityReturns = struct {
		result1 authorization.Identity
		result2 error
	}{result1, result2}
}

func (fake *IdentityProvider) GetIdentityReturnsOnCall(i int, result1 authorization.Identity, result2 error) {
	fake.getIdentityMutex.Lock()
	defer fake.getIdentityMutex.Unlock()
	fake.GetIdentityStub = nil
	if fake.getIdentityReturnsOnCall == nil {
		fake.getIdentityReturnsOnCall = make(map[int]struct {
			result1 authorization.Identity
			result2 error
		})
	}
	fake.getIdentityReturnsOnCall[i] = struct {
		result1 authorization.Identity
		result2 error
	}{result1, result2}
}

func (fake *IdentityProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getIdentityMutex.RLock()
	defer fake.getIdentityMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *IdentityProvider) recordInvocation(key string, args []interface{}) {
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

var _ apis.IdentityProvider = new(IdentityProvider)
