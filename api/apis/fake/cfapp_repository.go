// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/cf-k8s-controllers/api/apis"
	"code.cloudfoundry.org/cf-k8s-controllers/api/authorization"
	"code.cloudfoundry.org/cf-k8s-controllers/api/repositories"
)

type CFAppRepository struct {
	CreateAppStub        func(context.Context, authorization.Info, repositories.AppCreateMessage) (repositories.AppRecord, error)
	createAppMutex       sync.RWMutex
	createAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AppCreateMessage
	}
	createAppReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	createAppReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	CreateOrPatchAppEnvVarsStub        func(context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error)
	createOrPatchAppEnvVarsMutex       sync.RWMutex
	createOrPatchAppEnvVarsArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateOrPatchAppEnvVarsMessage
	}
	createOrPatchAppEnvVarsReturns struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}
	createOrPatchAppEnvVarsReturnsOnCall map[int]struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}
	FetchAppStub        func(context.Context, authorization.Info, string) (repositories.AppRecord, error)
	fetchAppMutex       sync.RWMutex
	fetchAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	fetchAppReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	fetchAppReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	FetchAppByNameAndSpaceStub        func(context.Context, authorization.Info, string, string) (repositories.AppRecord, error)
	fetchAppByNameAndSpaceMutex       sync.RWMutex
	fetchAppByNameAndSpaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}
	fetchAppByNameAndSpaceReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	fetchAppByNameAndSpaceReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	FetchAppListStub        func(context.Context, authorization.Info, repositories.AppListMessage) ([]repositories.AppRecord, error)
	fetchAppListMutex       sync.RWMutex
	fetchAppListArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AppListMessage
	}
	fetchAppListReturns struct {
		result1 []repositories.AppRecord
		result2 error
	}
	fetchAppListReturnsOnCall map[int]struct {
		result1 []repositories.AppRecord
		result2 error
	}
	FetchNamespaceStub        func(context.Context, authorization.Info, string) (repositories.SpaceRecord, error)
	fetchNamespaceMutex       sync.RWMutex
	fetchNamespaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	fetchNamespaceReturns struct {
		result1 repositories.SpaceRecord
		result2 error
	}
	fetchNamespaceReturnsOnCall map[int]struct {
		result1 repositories.SpaceRecord
		result2 error
	}
	SetAppDesiredStateStub        func(context.Context, authorization.Info, repositories.SetAppDesiredStateMessage) (repositories.AppRecord, error)
	setAppDesiredStateMutex       sync.RWMutex
	setAppDesiredStateArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.SetAppDesiredStateMessage
	}
	setAppDesiredStateReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	setAppDesiredStateReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	SetCurrentDropletStub        func(context.Context, authorization.Info, repositories.SetCurrentDropletMessage) (repositories.CurrentDropletRecord, error)
	setCurrentDropletMutex       sync.RWMutex
	setCurrentDropletArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.SetCurrentDropletMessage
	}
	setCurrentDropletReturns struct {
		result1 repositories.CurrentDropletRecord
		result2 error
	}
	setCurrentDropletReturnsOnCall map[int]struct {
		result1 repositories.CurrentDropletRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFAppRepository) CreateApp(arg1 context.Context, arg2 authorization.Info, arg3 repositories.AppCreateMessage) (repositories.AppRecord, error) {
	fake.createAppMutex.Lock()
	ret, specificReturn := fake.createAppReturnsOnCall[len(fake.createAppArgsForCall)]
	fake.createAppArgsForCall = append(fake.createAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AppCreateMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateAppStub
	fakeReturns := fake.createAppReturns
	fake.recordInvocation("CreateApp", []interface{}{arg1, arg2, arg3})
	fake.createAppMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) CreateAppCallCount() int {
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	return len(fake.createAppArgsForCall)
}

func (fake *CFAppRepository) CreateAppCalls(stub func(context.Context, authorization.Info, repositories.AppCreateMessage) (repositories.AppRecord, error)) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = stub
}

func (fake *CFAppRepository) CreateAppArgsForCall(i int) (context.Context, authorization.Info, repositories.AppCreateMessage) {
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	argsForCall := fake.createAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) CreateAppReturns(result1 repositories.AppRecord, result2 error) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = nil
	fake.createAppReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) CreateAppReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = nil
	if fake.createAppReturnsOnCall == nil {
		fake.createAppReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.createAppReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVars(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	ret, specificReturn := fake.createOrPatchAppEnvVarsReturnsOnCall[len(fake.createOrPatchAppEnvVarsArgsForCall)]
	fake.createOrPatchAppEnvVarsArgsForCall = append(fake.createOrPatchAppEnvVarsArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateOrPatchAppEnvVarsMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateOrPatchAppEnvVarsStub
	fakeReturns := fake.createOrPatchAppEnvVarsReturns
	fake.recordInvocation("CreateOrPatchAppEnvVars", []interface{}{arg1, arg2, arg3})
	fake.createOrPatchAppEnvVarsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsCallCount() int {
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	return len(fake.createOrPatchAppEnvVarsArgsForCall)
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsCalls(stub func(context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error)) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = stub
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) {
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	argsForCall := fake.createOrPatchAppEnvVarsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsReturns(result1 repositories.AppEnvVarsRecord, result2 error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = nil
	fake.createOrPatchAppEnvVarsReturns = struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) CreateOrPatchAppEnvVarsReturnsOnCall(i int, result1 repositories.AppEnvVarsRecord, result2 error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = nil
	if fake.createOrPatchAppEnvVarsReturnsOnCall == nil {
		fake.createOrPatchAppEnvVarsReturnsOnCall = make(map[int]struct {
			result1 repositories.AppEnvVarsRecord
			result2 error
		})
	}
	fake.createOrPatchAppEnvVarsReturnsOnCall[i] = struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchApp(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.AppRecord, error) {
	fake.fetchAppMutex.Lock()
	ret, specificReturn := fake.fetchAppReturnsOnCall[len(fake.fetchAppArgsForCall)]
	fake.fetchAppArgsForCall = append(fake.fetchAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.FetchAppStub
	fakeReturns := fake.fetchAppReturns
	fake.recordInvocation("FetchApp", []interface{}{arg1, arg2, arg3})
	fake.fetchAppMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) FetchAppCallCount() int {
	fake.fetchAppMutex.RLock()
	defer fake.fetchAppMutex.RUnlock()
	return len(fake.fetchAppArgsForCall)
}

func (fake *CFAppRepository) FetchAppCalls(stub func(context.Context, authorization.Info, string) (repositories.AppRecord, error)) {
	fake.fetchAppMutex.Lock()
	defer fake.fetchAppMutex.Unlock()
	fake.FetchAppStub = stub
}

func (fake *CFAppRepository) FetchAppArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.fetchAppMutex.RLock()
	defer fake.fetchAppMutex.RUnlock()
	argsForCall := fake.fetchAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) FetchAppReturns(result1 repositories.AppRecord, result2 error) {
	fake.fetchAppMutex.Lock()
	defer fake.fetchAppMutex.Unlock()
	fake.FetchAppStub = nil
	fake.fetchAppReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchAppReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.fetchAppMutex.Lock()
	defer fake.fetchAppMutex.Unlock()
	fake.FetchAppStub = nil
	if fake.fetchAppReturnsOnCall == nil {
		fake.fetchAppReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.fetchAppReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchAppByNameAndSpace(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string) (repositories.AppRecord, error) {
	fake.fetchAppByNameAndSpaceMutex.Lock()
	ret, specificReturn := fake.fetchAppByNameAndSpaceReturnsOnCall[len(fake.fetchAppByNameAndSpaceArgsForCall)]
	fake.fetchAppByNameAndSpaceArgsForCall = append(fake.fetchAppByNameAndSpaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.FetchAppByNameAndSpaceStub
	fakeReturns := fake.fetchAppByNameAndSpaceReturns
	fake.recordInvocation("FetchAppByNameAndSpace", []interface{}{arg1, arg2, arg3, arg4})
	fake.fetchAppByNameAndSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) FetchAppByNameAndSpaceCallCount() int {
	fake.fetchAppByNameAndSpaceMutex.RLock()
	defer fake.fetchAppByNameAndSpaceMutex.RUnlock()
	return len(fake.fetchAppByNameAndSpaceArgsForCall)
}

func (fake *CFAppRepository) FetchAppByNameAndSpaceCalls(stub func(context.Context, authorization.Info, string, string) (repositories.AppRecord, error)) {
	fake.fetchAppByNameAndSpaceMutex.Lock()
	defer fake.fetchAppByNameAndSpaceMutex.Unlock()
	fake.FetchAppByNameAndSpaceStub = stub
}

func (fake *CFAppRepository) FetchAppByNameAndSpaceArgsForCall(i int) (context.Context, authorization.Info, string, string) {
	fake.fetchAppByNameAndSpaceMutex.RLock()
	defer fake.fetchAppByNameAndSpaceMutex.RUnlock()
	argsForCall := fake.fetchAppByNameAndSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *CFAppRepository) FetchAppByNameAndSpaceReturns(result1 repositories.AppRecord, result2 error) {
	fake.fetchAppByNameAndSpaceMutex.Lock()
	defer fake.fetchAppByNameAndSpaceMutex.Unlock()
	fake.FetchAppByNameAndSpaceStub = nil
	fake.fetchAppByNameAndSpaceReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchAppByNameAndSpaceReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.fetchAppByNameAndSpaceMutex.Lock()
	defer fake.fetchAppByNameAndSpaceMutex.Unlock()
	fake.FetchAppByNameAndSpaceStub = nil
	if fake.fetchAppByNameAndSpaceReturnsOnCall == nil {
		fake.fetchAppByNameAndSpaceReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.fetchAppByNameAndSpaceReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchAppList(arg1 context.Context, arg2 authorization.Info, arg3 repositories.AppListMessage) ([]repositories.AppRecord, error) {
	fake.fetchAppListMutex.Lock()
	ret, specificReturn := fake.fetchAppListReturnsOnCall[len(fake.fetchAppListArgsForCall)]
	fake.fetchAppListArgsForCall = append(fake.fetchAppListArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AppListMessage
	}{arg1, arg2, arg3})
	stub := fake.FetchAppListStub
	fakeReturns := fake.fetchAppListReturns
	fake.recordInvocation("FetchAppList", []interface{}{arg1, arg2, arg3})
	fake.fetchAppListMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) FetchAppListCallCount() int {
	fake.fetchAppListMutex.RLock()
	defer fake.fetchAppListMutex.RUnlock()
	return len(fake.fetchAppListArgsForCall)
}

func (fake *CFAppRepository) FetchAppListCalls(stub func(context.Context, authorization.Info, repositories.AppListMessage) ([]repositories.AppRecord, error)) {
	fake.fetchAppListMutex.Lock()
	defer fake.fetchAppListMutex.Unlock()
	fake.FetchAppListStub = stub
}

func (fake *CFAppRepository) FetchAppListArgsForCall(i int) (context.Context, authorization.Info, repositories.AppListMessage) {
	fake.fetchAppListMutex.RLock()
	defer fake.fetchAppListMutex.RUnlock()
	argsForCall := fake.fetchAppListArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) FetchAppListReturns(result1 []repositories.AppRecord, result2 error) {
	fake.fetchAppListMutex.Lock()
	defer fake.fetchAppListMutex.Unlock()
	fake.FetchAppListStub = nil
	fake.fetchAppListReturns = struct {
		result1 []repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchAppListReturnsOnCall(i int, result1 []repositories.AppRecord, result2 error) {
	fake.fetchAppListMutex.Lock()
	defer fake.fetchAppListMutex.Unlock()
	fake.FetchAppListStub = nil
	if fake.fetchAppListReturnsOnCall == nil {
		fake.fetchAppListReturnsOnCall = make(map[int]struct {
			result1 []repositories.AppRecord
			result2 error
		})
	}
	fake.fetchAppListReturnsOnCall[i] = struct {
		result1 []repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchNamespace(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.SpaceRecord, error) {
	fake.fetchNamespaceMutex.Lock()
	ret, specificReturn := fake.fetchNamespaceReturnsOnCall[len(fake.fetchNamespaceArgsForCall)]
	fake.fetchNamespaceArgsForCall = append(fake.fetchNamespaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.FetchNamespaceStub
	fakeReturns := fake.fetchNamespaceReturns
	fake.recordInvocation("FetchNamespace", []interface{}{arg1, arg2, arg3})
	fake.fetchNamespaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) FetchNamespaceCallCount() int {
	fake.fetchNamespaceMutex.RLock()
	defer fake.fetchNamespaceMutex.RUnlock()
	return len(fake.fetchNamespaceArgsForCall)
}

func (fake *CFAppRepository) FetchNamespaceCalls(stub func(context.Context, authorization.Info, string) (repositories.SpaceRecord, error)) {
	fake.fetchNamespaceMutex.Lock()
	defer fake.fetchNamespaceMutex.Unlock()
	fake.FetchNamespaceStub = stub
}

func (fake *CFAppRepository) FetchNamespaceArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.fetchNamespaceMutex.RLock()
	defer fake.fetchNamespaceMutex.RUnlock()
	argsForCall := fake.fetchNamespaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) FetchNamespaceReturns(result1 repositories.SpaceRecord, result2 error) {
	fake.fetchNamespaceMutex.Lock()
	defer fake.fetchNamespaceMutex.Unlock()
	fake.FetchNamespaceStub = nil
	fake.fetchNamespaceReturns = struct {
		result1 repositories.SpaceRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) FetchNamespaceReturnsOnCall(i int, result1 repositories.SpaceRecord, result2 error) {
	fake.fetchNamespaceMutex.Lock()
	defer fake.fetchNamespaceMutex.Unlock()
	fake.FetchNamespaceStub = nil
	if fake.fetchNamespaceReturnsOnCall == nil {
		fake.fetchNamespaceReturnsOnCall = make(map[int]struct {
			result1 repositories.SpaceRecord
			result2 error
		})
	}
	fake.fetchNamespaceReturnsOnCall[i] = struct {
		result1 repositories.SpaceRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) SetAppDesiredState(arg1 context.Context, arg2 authorization.Info, arg3 repositories.SetAppDesiredStateMessage) (repositories.AppRecord, error) {
	fake.setAppDesiredStateMutex.Lock()
	ret, specificReturn := fake.setAppDesiredStateReturnsOnCall[len(fake.setAppDesiredStateArgsForCall)]
	fake.setAppDesiredStateArgsForCall = append(fake.setAppDesiredStateArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.SetAppDesiredStateMessage
	}{arg1, arg2, arg3})
	stub := fake.SetAppDesiredStateStub
	fakeReturns := fake.setAppDesiredStateReturns
	fake.recordInvocation("SetAppDesiredState", []interface{}{arg1, arg2, arg3})
	fake.setAppDesiredStateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) SetAppDesiredStateCallCount() int {
	fake.setAppDesiredStateMutex.RLock()
	defer fake.setAppDesiredStateMutex.RUnlock()
	return len(fake.setAppDesiredStateArgsForCall)
}

func (fake *CFAppRepository) SetAppDesiredStateCalls(stub func(context.Context, authorization.Info, repositories.SetAppDesiredStateMessage) (repositories.AppRecord, error)) {
	fake.setAppDesiredStateMutex.Lock()
	defer fake.setAppDesiredStateMutex.Unlock()
	fake.SetAppDesiredStateStub = stub
}

func (fake *CFAppRepository) SetAppDesiredStateArgsForCall(i int) (context.Context, authorization.Info, repositories.SetAppDesiredStateMessage) {
	fake.setAppDesiredStateMutex.RLock()
	defer fake.setAppDesiredStateMutex.RUnlock()
	argsForCall := fake.setAppDesiredStateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) SetAppDesiredStateReturns(result1 repositories.AppRecord, result2 error) {
	fake.setAppDesiredStateMutex.Lock()
	defer fake.setAppDesiredStateMutex.Unlock()
	fake.SetAppDesiredStateStub = nil
	fake.setAppDesiredStateReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) SetAppDesiredStateReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.setAppDesiredStateMutex.Lock()
	defer fake.setAppDesiredStateMutex.Unlock()
	fake.SetAppDesiredStateStub = nil
	if fake.setAppDesiredStateReturnsOnCall == nil {
		fake.setAppDesiredStateReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.setAppDesiredStateReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) SetCurrentDroplet(arg1 context.Context, arg2 authorization.Info, arg3 repositories.SetCurrentDropletMessage) (repositories.CurrentDropletRecord, error) {
	fake.setCurrentDropletMutex.Lock()
	ret, specificReturn := fake.setCurrentDropletReturnsOnCall[len(fake.setCurrentDropletArgsForCall)]
	fake.setCurrentDropletArgsForCall = append(fake.setCurrentDropletArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.SetCurrentDropletMessage
	}{arg1, arg2, arg3})
	stub := fake.SetCurrentDropletStub
	fakeReturns := fake.setCurrentDropletReturns
	fake.recordInvocation("SetCurrentDroplet", []interface{}{arg1, arg2, arg3})
	fake.setCurrentDropletMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFAppRepository) SetCurrentDropletCallCount() int {
	fake.setCurrentDropletMutex.RLock()
	defer fake.setCurrentDropletMutex.RUnlock()
	return len(fake.setCurrentDropletArgsForCall)
}

func (fake *CFAppRepository) SetCurrentDropletCalls(stub func(context.Context, authorization.Info, repositories.SetCurrentDropletMessage) (repositories.CurrentDropletRecord, error)) {
	fake.setCurrentDropletMutex.Lock()
	defer fake.setCurrentDropletMutex.Unlock()
	fake.SetCurrentDropletStub = stub
}

func (fake *CFAppRepository) SetCurrentDropletArgsForCall(i int) (context.Context, authorization.Info, repositories.SetCurrentDropletMessage) {
	fake.setCurrentDropletMutex.RLock()
	defer fake.setCurrentDropletMutex.RUnlock()
	argsForCall := fake.setCurrentDropletArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFAppRepository) SetCurrentDropletReturns(result1 repositories.CurrentDropletRecord, result2 error) {
	fake.setCurrentDropletMutex.Lock()
	defer fake.setCurrentDropletMutex.Unlock()
	fake.SetCurrentDropletStub = nil
	fake.setCurrentDropletReturns = struct {
		result1 repositories.CurrentDropletRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) SetCurrentDropletReturnsOnCall(i int, result1 repositories.CurrentDropletRecord, result2 error) {
	fake.setCurrentDropletMutex.Lock()
	defer fake.setCurrentDropletMutex.Unlock()
	fake.SetCurrentDropletStub = nil
	if fake.setCurrentDropletReturnsOnCall == nil {
		fake.setCurrentDropletReturnsOnCall = make(map[int]struct {
			result1 repositories.CurrentDropletRecord
			result2 error
		})
	}
	fake.setCurrentDropletReturnsOnCall[i] = struct {
		result1 repositories.CurrentDropletRecord
		result2 error
	}{result1, result2}
}

func (fake *CFAppRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	fake.fetchAppMutex.RLock()
	defer fake.fetchAppMutex.RUnlock()
	fake.fetchAppByNameAndSpaceMutex.RLock()
	defer fake.fetchAppByNameAndSpaceMutex.RUnlock()
	fake.fetchAppListMutex.RLock()
	defer fake.fetchAppListMutex.RUnlock()
	fake.fetchNamespaceMutex.RLock()
	defer fake.fetchNamespaceMutex.RUnlock()
	fake.setAppDesiredStateMutex.RLock()
	defer fake.setAppDesiredStateMutex.RUnlock()
	fake.setCurrentDropletMutex.RLock()
	defer fake.setCurrentDropletMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFAppRepository) recordInvocation(key string, args []interface{}) {
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

var _ apis.CFAppRepository = new(CFAppRepository)
