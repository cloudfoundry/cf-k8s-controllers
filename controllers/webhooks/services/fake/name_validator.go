// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"context"
	"sync"

	"code.cloudfoundry.org/korifi/controllers/webhooks"
	"code.cloudfoundry.org/korifi/controllers/webhooks/services"
	"github.com/go-logr/logr"
)

type NameValidator struct {
	ValidateCreateStub        func(context.Context, logr.Logger, string, string, string) *webhooks.ValidationError
	validateCreateMutex       sync.RWMutex
	validateCreateArgsForCall []struct {
		arg1 context.Context
		arg2 logr.Logger
		arg3 string
		arg4 string
		arg5 string
	}
	validateCreateReturns struct {
		result1 *webhooks.ValidationError
	}
	validateCreateReturnsOnCall map[int]struct {
		result1 *webhooks.ValidationError
	}
	ValidateDeleteStub        func(context.Context, logr.Logger, string, string) *webhooks.ValidationError
	validateDeleteMutex       sync.RWMutex
	validateDeleteArgsForCall []struct {
		arg1 context.Context
		arg2 logr.Logger
		arg3 string
		arg4 string
	}
	validateDeleteReturns struct {
		result1 *webhooks.ValidationError
	}
	validateDeleteReturnsOnCall map[int]struct {
		result1 *webhooks.ValidationError
	}
	ValidateUpdateStub        func(context.Context, logr.Logger, string, string, string, string) *webhooks.ValidationError
	validateUpdateMutex       sync.RWMutex
	validateUpdateArgsForCall []struct {
		arg1 context.Context
		arg2 logr.Logger
		arg3 string
		arg4 string
		arg5 string
		arg6 string
	}
	validateUpdateReturns struct {
		result1 *webhooks.ValidationError
	}
	validateUpdateReturnsOnCall map[int]struct {
		result1 *webhooks.ValidationError
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *NameValidator) ValidateCreate(arg1 context.Context, arg2 logr.Logger, arg3 string, arg4 string, arg5 string) *webhooks.ValidationError {
	fake.validateCreateMutex.Lock()
	ret, specificReturn := fake.validateCreateReturnsOnCall[len(fake.validateCreateArgsForCall)]
	fake.validateCreateArgsForCall = append(fake.validateCreateArgsForCall, struct {
		arg1 context.Context
		arg2 logr.Logger
		arg3 string
		arg4 string
		arg5 string
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.ValidateCreateStub
	fakeReturns := fake.validateCreateReturns
	fake.recordInvocation("ValidateCreate", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.validateCreateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *NameValidator) ValidateCreateCallCount() int {
	fake.validateCreateMutex.RLock()
	defer fake.validateCreateMutex.RUnlock()
	return len(fake.validateCreateArgsForCall)
}

func (fake *NameValidator) ValidateCreateCalls(stub func(context.Context, logr.Logger, string, string, string) *webhooks.ValidationError) {
	fake.validateCreateMutex.Lock()
	defer fake.validateCreateMutex.Unlock()
	fake.ValidateCreateStub = stub
}

func (fake *NameValidator) ValidateCreateArgsForCall(i int) (context.Context, logr.Logger, string, string, string) {
	fake.validateCreateMutex.RLock()
	defer fake.validateCreateMutex.RUnlock()
	argsForCall := fake.validateCreateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *NameValidator) ValidateCreateReturns(result1 *webhooks.ValidationError) {
	fake.validateCreateMutex.Lock()
	defer fake.validateCreateMutex.Unlock()
	fake.ValidateCreateStub = nil
	fake.validateCreateReturns = struct {
		result1 *webhooks.ValidationError
	}{result1}
}

func (fake *NameValidator) ValidateCreateReturnsOnCall(i int, result1 *webhooks.ValidationError) {
	fake.validateCreateMutex.Lock()
	defer fake.validateCreateMutex.Unlock()
	fake.ValidateCreateStub = nil
	if fake.validateCreateReturnsOnCall == nil {
		fake.validateCreateReturnsOnCall = make(map[int]struct {
			result1 *webhooks.ValidationError
		})
	}
	fake.validateCreateReturnsOnCall[i] = struct {
		result1 *webhooks.ValidationError
	}{result1}
}

func (fake *NameValidator) ValidateDelete(arg1 context.Context, arg2 logr.Logger, arg3 string, arg4 string) *webhooks.ValidationError {
	fake.validateDeleteMutex.Lock()
	ret, specificReturn := fake.validateDeleteReturnsOnCall[len(fake.validateDeleteArgsForCall)]
	fake.validateDeleteArgsForCall = append(fake.validateDeleteArgsForCall, struct {
		arg1 context.Context
		arg2 logr.Logger
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.ValidateDeleteStub
	fakeReturns := fake.validateDeleteReturns
	fake.recordInvocation("ValidateDelete", []interface{}{arg1, arg2, arg3, arg4})
	fake.validateDeleteMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *NameValidator) ValidateDeleteCallCount() int {
	fake.validateDeleteMutex.RLock()
	defer fake.validateDeleteMutex.RUnlock()
	return len(fake.validateDeleteArgsForCall)
}

func (fake *NameValidator) ValidateDeleteCalls(stub func(context.Context, logr.Logger, string, string) *webhooks.ValidationError) {
	fake.validateDeleteMutex.Lock()
	defer fake.validateDeleteMutex.Unlock()
	fake.ValidateDeleteStub = stub
}

func (fake *NameValidator) ValidateDeleteArgsForCall(i int) (context.Context, logr.Logger, string, string) {
	fake.validateDeleteMutex.RLock()
	defer fake.validateDeleteMutex.RUnlock()
	argsForCall := fake.validateDeleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *NameValidator) ValidateDeleteReturns(result1 *webhooks.ValidationError) {
	fake.validateDeleteMutex.Lock()
	defer fake.validateDeleteMutex.Unlock()
	fake.ValidateDeleteStub = nil
	fake.validateDeleteReturns = struct {
		result1 *webhooks.ValidationError
	}{result1}
}

func (fake *NameValidator) ValidateDeleteReturnsOnCall(i int, result1 *webhooks.ValidationError) {
	fake.validateDeleteMutex.Lock()
	defer fake.validateDeleteMutex.Unlock()
	fake.ValidateDeleteStub = nil
	if fake.validateDeleteReturnsOnCall == nil {
		fake.validateDeleteReturnsOnCall = make(map[int]struct {
			result1 *webhooks.ValidationError
		})
	}
	fake.validateDeleteReturnsOnCall[i] = struct {
		result1 *webhooks.ValidationError
	}{result1}
}

func (fake *NameValidator) ValidateUpdate(arg1 context.Context, arg2 logr.Logger, arg3 string, arg4 string, arg5 string, arg6 string) *webhooks.ValidationError {
	fake.validateUpdateMutex.Lock()
	ret, specificReturn := fake.validateUpdateReturnsOnCall[len(fake.validateUpdateArgsForCall)]
	fake.validateUpdateArgsForCall = append(fake.validateUpdateArgsForCall, struct {
		arg1 context.Context
		arg2 logr.Logger
		arg3 string
		arg4 string
		arg5 string
		arg6 string
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	stub := fake.ValidateUpdateStub
	fakeReturns := fake.validateUpdateReturns
	fake.recordInvocation("ValidateUpdate", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.validateUpdateMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *NameValidator) ValidateUpdateCallCount() int {
	fake.validateUpdateMutex.RLock()
	defer fake.validateUpdateMutex.RUnlock()
	return len(fake.validateUpdateArgsForCall)
}

func (fake *NameValidator) ValidateUpdateCalls(stub func(context.Context, logr.Logger, string, string, string, string) *webhooks.ValidationError) {
	fake.validateUpdateMutex.Lock()
	defer fake.validateUpdateMutex.Unlock()
	fake.ValidateUpdateStub = stub
}

func (fake *NameValidator) ValidateUpdateArgsForCall(i int) (context.Context, logr.Logger, string, string, string, string) {
	fake.validateUpdateMutex.RLock()
	defer fake.validateUpdateMutex.RUnlock()
	argsForCall := fake.validateUpdateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *NameValidator) ValidateUpdateReturns(result1 *webhooks.ValidationError) {
	fake.validateUpdateMutex.Lock()
	defer fake.validateUpdateMutex.Unlock()
	fake.ValidateUpdateStub = nil
	fake.validateUpdateReturns = struct {
		result1 *webhooks.ValidationError
	}{result1}
}

func (fake *NameValidator) ValidateUpdateReturnsOnCall(i int, result1 *webhooks.ValidationError) {
	fake.validateUpdateMutex.Lock()
	defer fake.validateUpdateMutex.Unlock()
	fake.ValidateUpdateStub = nil
	if fake.validateUpdateReturnsOnCall == nil {
		fake.validateUpdateReturnsOnCall = make(map[int]struct {
			result1 *webhooks.ValidationError
		})
	}
	fake.validateUpdateReturnsOnCall[i] = struct {
		result1 *webhooks.ValidationError
	}{result1}
}

func (fake *NameValidator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.validateCreateMutex.RLock()
	defer fake.validateCreateMutex.RUnlock()
	fake.validateDeleteMutex.RLock()
	defer fake.validateDeleteMutex.RUnlock()
	fake.validateUpdateMutex.RLock()
	defer fake.validateUpdateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *NameValidator) recordInvocation(key string, args []interface{}) {
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

var _ services.NameValidator = new(NameValidator)
