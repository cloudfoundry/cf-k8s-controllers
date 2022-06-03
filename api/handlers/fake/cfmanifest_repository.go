// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/api/authorization"
	"code.cloudfoundry.org/korifi/api/handlers"
	"code.cloudfoundry.org/korifi/api/repositories"
)

type CFManifestRepository struct {
	AddDestinationsToRouteStub        func(context.Context, authorization.Info, repositories.AddDestinationsToRouteMessage) (repositories.RouteRecord, error)
	addDestinationsToRouteMutex       sync.RWMutex
	addDestinationsToRouteArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AddDestinationsToRouteMessage
	}
	addDestinationsToRouteReturns struct {
		result1 repositories.RouteRecord
		result2 error
	}
	addDestinationsToRouteReturnsOnCall map[int]struct {
		result1 repositories.RouteRecord
		result2 error
	}
	CreateAppStub        func(context.Context, authorization.Info, repositories.CreateAppMessage) (repositories.AppRecord, error)
	createAppMutex       sync.RWMutex
	createAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateAppMessage
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
	CreateProcessStub        func(context.Context, authorization.Info, repositories.CreateProcessMessage) error
	createProcessMutex       sync.RWMutex
	createProcessArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateProcessMessage
	}
	createProcessReturns struct {
		result1 error
	}
	createProcessReturnsOnCall map[int]struct {
		result1 error
	}
	GetAppByNameAndSpaceStub        func(context.Context, authorization.Info, string, string) (repositories.AppRecord, error)
	getAppByNameAndSpaceMutex       sync.RWMutex
	getAppByNameAndSpaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}
	getAppByNameAndSpaceReturns struct {
		result1 repositories.AppRecord
		result2 error
	}
	getAppByNameAndSpaceReturnsOnCall map[int]struct {
		result1 repositories.AppRecord
		result2 error
	}
	GetDomainByNameStub        func(context.Context, authorization.Info, string) (repositories.DomainRecord, error)
	getDomainByNameMutex       sync.RWMutex
	getDomainByNameArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getDomainByNameReturns struct {
		result1 repositories.DomainRecord
		result2 error
	}
	getDomainByNameReturnsOnCall map[int]struct {
		result1 repositories.DomainRecord
		result2 error
	}
	GetOrCreateRouteStub        func(context.Context, authorization.Info, repositories.CreateRouteMessage) (repositories.RouteRecord, error)
	getOrCreateRouteMutex       sync.RWMutex
	getOrCreateRouteArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateRouteMessage
	}
	getOrCreateRouteReturns struct {
		result1 repositories.RouteRecord
		result2 error
	}
	getOrCreateRouteReturnsOnCall map[int]struct {
		result1 repositories.RouteRecord
		result2 error
	}
	GetProcessByAppTypeAndSpaceStub        func(context.Context, authorization.Info, string, string, string) (repositories.ProcessRecord, error)
	getProcessByAppTypeAndSpaceMutex       sync.RWMutex
	getProcessByAppTypeAndSpaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
		arg5 string
	}
	getProcessByAppTypeAndSpaceReturns struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	getProcessByAppTypeAndSpaceReturnsOnCall map[int]struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	GetSpaceStub        func(context.Context, authorization.Info, string) (repositories.SpaceRecord, error)
	getSpaceMutex       sync.RWMutex
	getSpaceArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}
	getSpaceReturns struct {
		result1 repositories.SpaceRecord
		result2 error
	}
	getSpaceReturnsOnCall map[int]struct {
		result1 repositories.SpaceRecord
		result2 error
	}
	ListRoutesForAppStub        func(context.Context, authorization.Info, string, string) ([]repositories.RouteRecord, error)
	listRoutesForAppMutex       sync.RWMutex
	listRoutesForAppArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}
	listRoutesForAppReturns struct {
		result1 []repositories.RouteRecord
		result2 error
	}
	listRoutesForAppReturnsOnCall map[int]struct {
		result1 []repositories.RouteRecord
		result2 error
	}
	PatchProcessStub        func(context.Context, authorization.Info, repositories.PatchProcessMessage) (repositories.ProcessRecord, error)
	patchProcessMutex       sync.RWMutex
	patchProcessArgsForCall []struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.PatchProcessMessage
	}
	patchProcessReturns struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	patchProcessReturnsOnCall map[int]struct {
		result1 repositories.ProcessRecord
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CFManifestRepository) AddDestinationsToRoute(arg1 context.Context, arg2 authorization.Info, arg3 repositories.AddDestinationsToRouteMessage) (repositories.RouteRecord, error) {
	fake.addDestinationsToRouteMutex.Lock()
	ret, specificReturn := fake.addDestinationsToRouteReturnsOnCall[len(fake.addDestinationsToRouteArgsForCall)]
	fake.addDestinationsToRouteArgsForCall = append(fake.addDestinationsToRouteArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.AddDestinationsToRouteMessage
	}{arg1, arg2, arg3})
	stub := fake.AddDestinationsToRouteStub
	fakeReturns := fake.addDestinationsToRouteReturns
	fake.recordInvocation("AddDestinationsToRoute", []interface{}{arg1, arg2, arg3})
	fake.addDestinationsToRouteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) AddDestinationsToRouteCallCount() int {
	fake.addDestinationsToRouteMutex.RLock()
	defer fake.addDestinationsToRouteMutex.RUnlock()
	return len(fake.addDestinationsToRouteArgsForCall)
}

func (fake *CFManifestRepository) AddDestinationsToRouteCalls(stub func(context.Context, authorization.Info, repositories.AddDestinationsToRouteMessage) (repositories.RouteRecord, error)) {
	fake.addDestinationsToRouteMutex.Lock()
	defer fake.addDestinationsToRouteMutex.Unlock()
	fake.AddDestinationsToRouteStub = stub
}

func (fake *CFManifestRepository) AddDestinationsToRouteArgsForCall(i int) (context.Context, authorization.Info, repositories.AddDestinationsToRouteMessage) {
	fake.addDestinationsToRouteMutex.RLock()
	defer fake.addDestinationsToRouteMutex.RUnlock()
	argsForCall := fake.addDestinationsToRouteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) AddDestinationsToRouteReturns(result1 repositories.RouteRecord, result2 error) {
	fake.addDestinationsToRouteMutex.Lock()
	defer fake.addDestinationsToRouteMutex.Unlock()
	fake.AddDestinationsToRouteStub = nil
	fake.addDestinationsToRouteReturns = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) AddDestinationsToRouteReturnsOnCall(i int, result1 repositories.RouteRecord, result2 error) {
	fake.addDestinationsToRouteMutex.Lock()
	defer fake.addDestinationsToRouteMutex.Unlock()
	fake.AddDestinationsToRouteStub = nil
	if fake.addDestinationsToRouteReturnsOnCall == nil {
		fake.addDestinationsToRouteReturnsOnCall = make(map[int]struct {
			result1 repositories.RouteRecord
			result2 error
		})
	}
	fake.addDestinationsToRouteReturnsOnCall[i] = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) CreateApp(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateAppMessage) (repositories.AppRecord, error) {
	fake.createAppMutex.Lock()
	ret, specificReturn := fake.createAppReturnsOnCall[len(fake.createAppArgsForCall)]
	fake.createAppArgsForCall = append(fake.createAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateAppMessage
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

func (fake *CFManifestRepository) CreateAppCallCount() int {
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	return len(fake.createAppArgsForCall)
}

func (fake *CFManifestRepository) CreateAppCalls(stub func(context.Context, authorization.Info, repositories.CreateAppMessage) (repositories.AppRecord, error)) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = stub
}

func (fake *CFManifestRepository) CreateAppArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateAppMessage) {
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	argsForCall := fake.createAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) CreateAppReturns(result1 repositories.AppRecord, result2 error) {
	fake.createAppMutex.Lock()
	defer fake.createAppMutex.Unlock()
	fake.CreateAppStub = nil
	fake.createAppReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) CreateAppReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
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

func (fake *CFManifestRepository) CreateOrPatchAppEnvVars(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error) {
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

func (fake *CFManifestRepository) CreateOrPatchAppEnvVarsCallCount() int {
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	return len(fake.createOrPatchAppEnvVarsArgsForCall)
}

func (fake *CFManifestRepository) CreateOrPatchAppEnvVarsCalls(stub func(context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) (repositories.AppEnvVarsRecord, error)) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = stub
}

func (fake *CFManifestRepository) CreateOrPatchAppEnvVarsArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateOrPatchAppEnvVarsMessage) {
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	argsForCall := fake.createOrPatchAppEnvVarsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) CreateOrPatchAppEnvVarsReturns(result1 repositories.AppEnvVarsRecord, result2 error) {
	fake.createOrPatchAppEnvVarsMutex.Lock()
	defer fake.createOrPatchAppEnvVarsMutex.Unlock()
	fake.CreateOrPatchAppEnvVarsStub = nil
	fake.createOrPatchAppEnvVarsReturns = struct {
		result1 repositories.AppEnvVarsRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) CreateOrPatchAppEnvVarsReturnsOnCall(i int, result1 repositories.AppEnvVarsRecord, result2 error) {
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

func (fake *CFManifestRepository) CreateProcess(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateProcessMessage) error {
	fake.createProcessMutex.Lock()
	ret, specificReturn := fake.createProcessReturnsOnCall[len(fake.createProcessArgsForCall)]
	fake.createProcessArgsForCall = append(fake.createProcessArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateProcessMessage
	}{arg1, arg2, arg3})
	stub := fake.CreateProcessStub
	fakeReturns := fake.createProcessReturns
	fake.recordInvocation("CreateProcess", []interface{}{arg1, arg2, arg3})
	fake.createProcessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CFManifestRepository) CreateProcessCallCount() int {
	fake.createProcessMutex.RLock()
	defer fake.createProcessMutex.RUnlock()
	return len(fake.createProcessArgsForCall)
}

func (fake *CFManifestRepository) CreateProcessCalls(stub func(context.Context, authorization.Info, repositories.CreateProcessMessage) error) {
	fake.createProcessMutex.Lock()
	defer fake.createProcessMutex.Unlock()
	fake.CreateProcessStub = stub
}

func (fake *CFManifestRepository) CreateProcessArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateProcessMessage) {
	fake.createProcessMutex.RLock()
	defer fake.createProcessMutex.RUnlock()
	argsForCall := fake.createProcessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) CreateProcessReturns(result1 error) {
	fake.createProcessMutex.Lock()
	defer fake.createProcessMutex.Unlock()
	fake.CreateProcessStub = nil
	fake.createProcessReturns = struct {
		result1 error
	}{result1}
}

func (fake *CFManifestRepository) CreateProcessReturnsOnCall(i int, result1 error) {
	fake.createProcessMutex.Lock()
	defer fake.createProcessMutex.Unlock()
	fake.CreateProcessStub = nil
	if fake.createProcessReturnsOnCall == nil {
		fake.createProcessReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createProcessReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CFManifestRepository) GetAppByNameAndSpace(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string) (repositories.AppRecord, error) {
	fake.getAppByNameAndSpaceMutex.Lock()
	ret, specificReturn := fake.getAppByNameAndSpaceReturnsOnCall[len(fake.getAppByNameAndSpaceArgsForCall)]
	fake.getAppByNameAndSpaceArgsForCall = append(fake.getAppByNameAndSpaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.GetAppByNameAndSpaceStub
	fakeReturns := fake.getAppByNameAndSpaceReturns
	fake.recordInvocation("GetAppByNameAndSpace", []interface{}{arg1, arg2, arg3, arg4})
	fake.getAppByNameAndSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) GetAppByNameAndSpaceCallCount() int {
	fake.getAppByNameAndSpaceMutex.RLock()
	defer fake.getAppByNameAndSpaceMutex.RUnlock()
	return len(fake.getAppByNameAndSpaceArgsForCall)
}

func (fake *CFManifestRepository) GetAppByNameAndSpaceCalls(stub func(context.Context, authorization.Info, string, string) (repositories.AppRecord, error)) {
	fake.getAppByNameAndSpaceMutex.Lock()
	defer fake.getAppByNameAndSpaceMutex.Unlock()
	fake.GetAppByNameAndSpaceStub = stub
}

func (fake *CFManifestRepository) GetAppByNameAndSpaceArgsForCall(i int) (context.Context, authorization.Info, string, string) {
	fake.getAppByNameAndSpaceMutex.RLock()
	defer fake.getAppByNameAndSpaceMutex.RUnlock()
	argsForCall := fake.getAppByNameAndSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *CFManifestRepository) GetAppByNameAndSpaceReturns(result1 repositories.AppRecord, result2 error) {
	fake.getAppByNameAndSpaceMutex.Lock()
	defer fake.getAppByNameAndSpaceMutex.Unlock()
	fake.GetAppByNameAndSpaceStub = nil
	fake.getAppByNameAndSpaceReturns = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetAppByNameAndSpaceReturnsOnCall(i int, result1 repositories.AppRecord, result2 error) {
	fake.getAppByNameAndSpaceMutex.Lock()
	defer fake.getAppByNameAndSpaceMutex.Unlock()
	fake.GetAppByNameAndSpaceStub = nil
	if fake.getAppByNameAndSpaceReturnsOnCall == nil {
		fake.getAppByNameAndSpaceReturnsOnCall = make(map[int]struct {
			result1 repositories.AppRecord
			result2 error
		})
	}
	fake.getAppByNameAndSpaceReturnsOnCall[i] = struct {
		result1 repositories.AppRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetDomainByName(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.DomainRecord, error) {
	fake.getDomainByNameMutex.Lock()
	ret, specificReturn := fake.getDomainByNameReturnsOnCall[len(fake.getDomainByNameArgsForCall)]
	fake.getDomainByNameArgsForCall = append(fake.getDomainByNameArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetDomainByNameStub
	fakeReturns := fake.getDomainByNameReturns
	fake.recordInvocation("GetDomainByName", []interface{}{arg1, arg2, arg3})
	fake.getDomainByNameMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) GetDomainByNameCallCount() int {
	fake.getDomainByNameMutex.RLock()
	defer fake.getDomainByNameMutex.RUnlock()
	return len(fake.getDomainByNameArgsForCall)
}

func (fake *CFManifestRepository) GetDomainByNameCalls(stub func(context.Context, authorization.Info, string) (repositories.DomainRecord, error)) {
	fake.getDomainByNameMutex.Lock()
	defer fake.getDomainByNameMutex.Unlock()
	fake.GetDomainByNameStub = stub
}

func (fake *CFManifestRepository) GetDomainByNameArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getDomainByNameMutex.RLock()
	defer fake.getDomainByNameMutex.RUnlock()
	argsForCall := fake.getDomainByNameArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) GetDomainByNameReturns(result1 repositories.DomainRecord, result2 error) {
	fake.getDomainByNameMutex.Lock()
	defer fake.getDomainByNameMutex.Unlock()
	fake.GetDomainByNameStub = nil
	fake.getDomainByNameReturns = struct {
		result1 repositories.DomainRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetDomainByNameReturnsOnCall(i int, result1 repositories.DomainRecord, result2 error) {
	fake.getDomainByNameMutex.Lock()
	defer fake.getDomainByNameMutex.Unlock()
	fake.GetDomainByNameStub = nil
	if fake.getDomainByNameReturnsOnCall == nil {
		fake.getDomainByNameReturnsOnCall = make(map[int]struct {
			result1 repositories.DomainRecord
			result2 error
		})
	}
	fake.getDomainByNameReturnsOnCall[i] = struct {
		result1 repositories.DomainRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetOrCreateRoute(arg1 context.Context, arg2 authorization.Info, arg3 repositories.CreateRouteMessage) (repositories.RouteRecord, error) {
	fake.getOrCreateRouteMutex.Lock()
	ret, specificReturn := fake.getOrCreateRouteReturnsOnCall[len(fake.getOrCreateRouteArgsForCall)]
	fake.getOrCreateRouteArgsForCall = append(fake.getOrCreateRouteArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.CreateRouteMessage
	}{arg1, arg2, arg3})
	stub := fake.GetOrCreateRouteStub
	fakeReturns := fake.getOrCreateRouteReturns
	fake.recordInvocation("GetOrCreateRoute", []interface{}{arg1, arg2, arg3})
	fake.getOrCreateRouteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) GetOrCreateRouteCallCount() int {
	fake.getOrCreateRouteMutex.RLock()
	defer fake.getOrCreateRouteMutex.RUnlock()
	return len(fake.getOrCreateRouteArgsForCall)
}

func (fake *CFManifestRepository) GetOrCreateRouteCalls(stub func(context.Context, authorization.Info, repositories.CreateRouteMessage) (repositories.RouteRecord, error)) {
	fake.getOrCreateRouteMutex.Lock()
	defer fake.getOrCreateRouteMutex.Unlock()
	fake.GetOrCreateRouteStub = stub
}

func (fake *CFManifestRepository) GetOrCreateRouteArgsForCall(i int) (context.Context, authorization.Info, repositories.CreateRouteMessage) {
	fake.getOrCreateRouteMutex.RLock()
	defer fake.getOrCreateRouteMutex.RUnlock()
	argsForCall := fake.getOrCreateRouteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) GetOrCreateRouteReturns(result1 repositories.RouteRecord, result2 error) {
	fake.getOrCreateRouteMutex.Lock()
	defer fake.getOrCreateRouteMutex.Unlock()
	fake.GetOrCreateRouteStub = nil
	fake.getOrCreateRouteReturns = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetOrCreateRouteReturnsOnCall(i int, result1 repositories.RouteRecord, result2 error) {
	fake.getOrCreateRouteMutex.Lock()
	defer fake.getOrCreateRouteMutex.Unlock()
	fake.GetOrCreateRouteStub = nil
	if fake.getOrCreateRouteReturnsOnCall == nil {
		fake.getOrCreateRouteReturnsOnCall = make(map[int]struct {
			result1 repositories.RouteRecord
			result2 error
		})
	}
	fake.getOrCreateRouteReturnsOnCall[i] = struct {
		result1 repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetProcessByAppTypeAndSpace(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string, arg5 string) (repositories.ProcessRecord, error) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	ret, specificReturn := fake.getProcessByAppTypeAndSpaceReturnsOnCall[len(fake.getProcessByAppTypeAndSpaceArgsForCall)]
	fake.getProcessByAppTypeAndSpaceArgsForCall = append(fake.getProcessByAppTypeAndSpaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.GetProcessByAppTypeAndSpaceStub
	fakeReturns := fake.getProcessByAppTypeAndSpaceReturns
	fake.recordInvocation("GetProcessByAppTypeAndSpace", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) GetProcessByAppTypeAndSpaceCallCount() int {
	fake.getProcessByAppTypeAndSpaceMutex.RLock()
	defer fake.getProcessByAppTypeAndSpaceMutex.RUnlock()
	return len(fake.getProcessByAppTypeAndSpaceArgsForCall)
}

func (fake *CFManifestRepository) GetProcessByAppTypeAndSpaceCalls(stub func(context.Context, authorization.Info, string, string, string) (repositories.ProcessRecord, error)) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	defer fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	fake.GetProcessByAppTypeAndSpaceStub = stub
}

func (fake *CFManifestRepository) GetProcessByAppTypeAndSpaceArgsForCall(i int) (context.Context, authorization.Info, string, string, string) {
	fake.getProcessByAppTypeAndSpaceMutex.RLock()
	defer fake.getProcessByAppTypeAndSpaceMutex.RUnlock()
	argsForCall := fake.getProcessByAppTypeAndSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *CFManifestRepository) GetProcessByAppTypeAndSpaceReturns(result1 repositories.ProcessRecord, result2 error) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	defer fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	fake.GetProcessByAppTypeAndSpaceStub = nil
	fake.getProcessByAppTypeAndSpaceReturns = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetProcessByAppTypeAndSpaceReturnsOnCall(i int, result1 repositories.ProcessRecord, result2 error) {
	fake.getProcessByAppTypeAndSpaceMutex.Lock()
	defer fake.getProcessByAppTypeAndSpaceMutex.Unlock()
	fake.GetProcessByAppTypeAndSpaceStub = nil
	if fake.getProcessByAppTypeAndSpaceReturnsOnCall == nil {
		fake.getProcessByAppTypeAndSpaceReturnsOnCall = make(map[int]struct {
			result1 repositories.ProcessRecord
			result2 error
		})
	}
	fake.getProcessByAppTypeAndSpaceReturnsOnCall[i] = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetSpace(arg1 context.Context, arg2 authorization.Info, arg3 string) (repositories.SpaceRecord, error) {
	fake.getSpaceMutex.Lock()
	ret, specificReturn := fake.getSpaceReturnsOnCall[len(fake.getSpaceArgsForCall)]
	fake.getSpaceArgsForCall = append(fake.getSpaceArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetSpaceStub
	fakeReturns := fake.getSpaceReturns
	fake.recordInvocation("GetSpace", []interface{}{arg1, arg2, arg3})
	fake.getSpaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) GetSpaceCallCount() int {
	fake.getSpaceMutex.RLock()
	defer fake.getSpaceMutex.RUnlock()
	return len(fake.getSpaceArgsForCall)
}

func (fake *CFManifestRepository) GetSpaceCalls(stub func(context.Context, authorization.Info, string) (repositories.SpaceRecord, error)) {
	fake.getSpaceMutex.Lock()
	defer fake.getSpaceMutex.Unlock()
	fake.GetSpaceStub = stub
}

func (fake *CFManifestRepository) GetSpaceArgsForCall(i int) (context.Context, authorization.Info, string) {
	fake.getSpaceMutex.RLock()
	defer fake.getSpaceMutex.RUnlock()
	argsForCall := fake.getSpaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) GetSpaceReturns(result1 repositories.SpaceRecord, result2 error) {
	fake.getSpaceMutex.Lock()
	defer fake.getSpaceMutex.Unlock()
	fake.GetSpaceStub = nil
	fake.getSpaceReturns = struct {
		result1 repositories.SpaceRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) GetSpaceReturnsOnCall(i int, result1 repositories.SpaceRecord, result2 error) {
	fake.getSpaceMutex.Lock()
	defer fake.getSpaceMutex.Unlock()
	fake.GetSpaceStub = nil
	if fake.getSpaceReturnsOnCall == nil {
		fake.getSpaceReturnsOnCall = make(map[int]struct {
			result1 repositories.SpaceRecord
			result2 error
		})
	}
	fake.getSpaceReturnsOnCall[i] = struct {
		result1 repositories.SpaceRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) ListRoutesForApp(arg1 context.Context, arg2 authorization.Info, arg3 string, arg4 string) ([]repositories.RouteRecord, error) {
	fake.listRoutesForAppMutex.Lock()
	ret, specificReturn := fake.listRoutesForAppReturnsOnCall[len(fake.listRoutesForAppArgsForCall)]
	fake.listRoutesForAppArgsForCall = append(fake.listRoutesForAppArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.ListRoutesForAppStub
	fakeReturns := fake.listRoutesForAppReturns
	fake.recordInvocation("ListRoutesForApp", []interface{}{arg1, arg2, arg3, arg4})
	fake.listRoutesForAppMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) ListRoutesForAppCallCount() int {
	fake.listRoutesForAppMutex.RLock()
	defer fake.listRoutesForAppMutex.RUnlock()
	return len(fake.listRoutesForAppArgsForCall)
}

func (fake *CFManifestRepository) ListRoutesForAppCalls(stub func(context.Context, authorization.Info, string, string) ([]repositories.RouteRecord, error)) {
	fake.listRoutesForAppMutex.Lock()
	defer fake.listRoutesForAppMutex.Unlock()
	fake.ListRoutesForAppStub = stub
}

func (fake *CFManifestRepository) ListRoutesForAppArgsForCall(i int) (context.Context, authorization.Info, string, string) {
	fake.listRoutesForAppMutex.RLock()
	defer fake.listRoutesForAppMutex.RUnlock()
	argsForCall := fake.listRoutesForAppArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *CFManifestRepository) ListRoutesForAppReturns(result1 []repositories.RouteRecord, result2 error) {
	fake.listRoutesForAppMutex.Lock()
	defer fake.listRoutesForAppMutex.Unlock()
	fake.ListRoutesForAppStub = nil
	fake.listRoutesForAppReturns = struct {
		result1 []repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) ListRoutesForAppReturnsOnCall(i int, result1 []repositories.RouteRecord, result2 error) {
	fake.listRoutesForAppMutex.Lock()
	defer fake.listRoutesForAppMutex.Unlock()
	fake.ListRoutesForAppStub = nil
	if fake.listRoutesForAppReturnsOnCall == nil {
		fake.listRoutesForAppReturnsOnCall = make(map[int]struct {
			result1 []repositories.RouteRecord
			result2 error
		})
	}
	fake.listRoutesForAppReturnsOnCall[i] = struct {
		result1 []repositories.RouteRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) PatchProcess(arg1 context.Context, arg2 authorization.Info, arg3 repositories.PatchProcessMessage) (repositories.ProcessRecord, error) {
	fake.patchProcessMutex.Lock()
	ret, specificReturn := fake.patchProcessReturnsOnCall[len(fake.patchProcessArgsForCall)]
	fake.patchProcessArgsForCall = append(fake.patchProcessArgsForCall, struct {
		arg1 context.Context
		arg2 authorization.Info
		arg3 repositories.PatchProcessMessage
	}{arg1, arg2, arg3})
	stub := fake.PatchProcessStub
	fakeReturns := fake.patchProcessReturns
	fake.recordInvocation("PatchProcess", []interface{}{arg1, arg2, arg3})
	fake.patchProcessMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CFManifestRepository) PatchProcessCallCount() int {
	fake.patchProcessMutex.RLock()
	defer fake.patchProcessMutex.RUnlock()
	return len(fake.patchProcessArgsForCall)
}

func (fake *CFManifestRepository) PatchProcessCalls(stub func(context.Context, authorization.Info, repositories.PatchProcessMessage) (repositories.ProcessRecord, error)) {
	fake.patchProcessMutex.Lock()
	defer fake.patchProcessMutex.Unlock()
	fake.PatchProcessStub = stub
}

func (fake *CFManifestRepository) PatchProcessArgsForCall(i int) (context.Context, authorization.Info, repositories.PatchProcessMessage) {
	fake.patchProcessMutex.RLock()
	defer fake.patchProcessMutex.RUnlock()
	argsForCall := fake.patchProcessArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CFManifestRepository) PatchProcessReturns(result1 repositories.ProcessRecord, result2 error) {
	fake.patchProcessMutex.Lock()
	defer fake.patchProcessMutex.Unlock()
	fake.PatchProcessStub = nil
	fake.patchProcessReturns = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) PatchProcessReturnsOnCall(i int, result1 repositories.ProcessRecord, result2 error) {
	fake.patchProcessMutex.Lock()
	defer fake.patchProcessMutex.Unlock()
	fake.PatchProcessStub = nil
	if fake.patchProcessReturnsOnCall == nil {
		fake.patchProcessReturnsOnCall = make(map[int]struct {
			result1 repositories.ProcessRecord
			result2 error
		})
	}
	fake.patchProcessReturnsOnCall[i] = struct {
		result1 repositories.ProcessRecord
		result2 error
	}{result1, result2}
}

func (fake *CFManifestRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addDestinationsToRouteMutex.RLock()
	defer fake.addDestinationsToRouteMutex.RUnlock()
	fake.createAppMutex.RLock()
	defer fake.createAppMutex.RUnlock()
	fake.createOrPatchAppEnvVarsMutex.RLock()
	defer fake.createOrPatchAppEnvVarsMutex.RUnlock()
	fake.createProcessMutex.RLock()
	defer fake.createProcessMutex.RUnlock()
	fake.getAppByNameAndSpaceMutex.RLock()
	defer fake.getAppByNameAndSpaceMutex.RUnlock()
	fake.getDomainByNameMutex.RLock()
	defer fake.getDomainByNameMutex.RUnlock()
	fake.getOrCreateRouteMutex.RLock()
	defer fake.getOrCreateRouteMutex.RUnlock()
	fake.getProcessByAppTypeAndSpaceMutex.RLock()
	defer fake.getProcessByAppTypeAndSpaceMutex.RUnlock()
	fake.getSpaceMutex.RLock()
	defer fake.getSpaceMutex.RUnlock()
	fake.listRoutesForAppMutex.RLock()
	defer fake.listRoutesForAppMutex.RUnlock()
	fake.patchProcessMutex.RLock()
	defer fake.patchProcessMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CFManifestRepository) recordInvocation(key string, args []interface{}) {
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

var _ handlers.CFManifestRepository = new(CFManifestRepository)
