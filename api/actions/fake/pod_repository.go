// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/actions"
	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type PodRepository struct {
	ListPodStatsStub        func(context.Context, authorization.Info, repositories.ListPodStatsMessage) ([]repositories.PodStatsRecord, error)
	listPodStatsMutex       sync.RWMutex
	listPodStatsArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListPodStatsMessage
	}
	listPodStatsReturns struct {
		result1 []repositories.PodStatsRecord
		result2 error
	}
	listPodStatsReturnsOnCall map[int]struct {
		result1 []repositories.PodStatsRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *PodRepository) ListPodStats(arg1 context.Context, arg2 authorization.Info, arg3 repositories.ListPodStatsMessage) ([]repositories.PodStatsRecord, error) {
	fake.listPodStatsMutex.Lock()
	ret, specificReturn := fake.listPodStatsReturnsOnCall[len(fake.listPodStatsArgsForCall)]
	fake.listPodStatsArgsForCall = append(fake.listPodStatsArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.ListPodStatsMessage
	}{arg1, arg2, arg3})
	stub := fake.ListPodStatsStub
	fakeReturns := fake.listPodStatsReturns
	fake.recordInvocation("ListPodStats", []interface{}{arg1, arg2, arg3})
	fake.listPodStatsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *PodRepository) ListPodStatsCallCount() int {
	fake.listPodStatsMutex.RLock()
	defer fake.listPodStatsMutex.RUnlock()
	return len(fake.listPodStatsArgsForCall)
}

func (fake *PodRepository) ListPodStatsCalls(stub func(context.Context, authorization.Info, repositories.ListPodStatsMessage) ([]repositories.PodStatsRecord, error)) {
	fake.listPodStatsMutex.Lock()
	defer fake.listPodStatsMutex.Unlock()
	fake.ListPodStatsStub = stub
}

func (fake *PodRepository) ListPodStatsArgsForCall(i int) (context.Context, authorization.Info, repositories.ListPodStatsMessage) {
	fake.listPodStatsMutex.RLock()
	defer fake.listPodStatsMutex.RUnlock()
	argsForCall := fake.listPodStatsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *PodRepository) ListPodStatsReturns(result1 []repositories.PodStatsRecord, result2 error) {
	fake.listPodStatsMutex.Lock()
	defer fake.listPodStatsMutex.Unlock()
	fake.ListPodStatsStub = nil
	fake.listPodStatsReturns = struct {
		result1 []repositories.PodStatsRecord
		result2 error
	}{result1, result2}
}

func (fake *PodRepository) ListPodStatsReturnsOnCall(i int, result1 []repositories.PodStatsRecord, result2 error) {
	fake.listPodStatsMutex.Lock()
	defer fake.listPodStatsMutex.Unlock()
	fake.ListPodStatsStub = nil
	if fake.listPodStatsReturnsOnCall == nil {
		fake.listPodStatsReturnsOnCall = make(map[int]struct {
			result1 []repositories.PodStatsRecord
			result2 error
		})
	}
	fake.listPodStatsReturnsOnCall[i] = struct {
		result1 []repositories.PodStatsRecord
		result2 error
	}{result1, result2}
}

func (fake *PodRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listPodStatsMutex.RLock()
	defer fake.listPodStatsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *PodRepository) recordInvocation(key string, args []interface{}) {
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

var _ actions.PodRepository = new(PodRepository)
