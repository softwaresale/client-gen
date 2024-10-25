// Code generated by mockery v2.46.3. DO NOT EDIT.

package outputs_mocks

import (
	outputs "github.com/softwaresale/client-gen/v2/internal/codegen/outputs"
	types "github.com/softwaresale/client-gen/v2/internal/types"
	mock "github.com/stretchr/testify/mock"
)

// MockCompilerOutputsManager is an autogenerated mock type for the CompilerOutputsManager type
type MockCompilerOutputsManager struct {
	mock.Mock
}

// ComputeModelLocation provides a mock function with given fields: model
func (_m *MockCompilerOutputsManager) ComputeModelLocation(model types.EntitySpec) (outputs.CompilerOutputLocation, error) {
	ret := _m.Called(model)

	if len(ret) == 0 {
		panic("no return value specified for ComputeModelLocation")
	}

	var r0 outputs.CompilerOutputLocation
	var r1 error
	if rf, ok := ret.Get(0).(func(types.EntitySpec) (outputs.CompilerOutputLocation, error)); ok {
		return rf(model)
	}
	if rf, ok := ret.Get(0).(func(types.EntitySpec) outputs.CompilerOutputLocation); ok {
		r0 = rf(model)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(outputs.CompilerOutputLocation)
		}
	}

	if rf, ok := ret.Get(1).(func(types.EntitySpec) error); ok {
		r1 = rf(model)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ComputeServiceLocation provides a mock function with given fields: serviceDef
func (_m *MockCompilerOutputsManager) ComputeServiceLocation(serviceDef types.ServiceDefinition) (outputs.CompilerOutputLocation, error) {
	ret := _m.Called(serviceDef)

	if len(ret) == 0 {
		panic("no return value specified for ComputeServiceLocation")
	}

	var r0 outputs.CompilerOutputLocation
	var r1 error
	if rf, ok := ret.Get(0).(func(types.ServiceDefinition) (outputs.CompilerOutputLocation, error)); ok {
		return rf(serviceDef)
	}
	if rf, ok := ret.Get(0).(func(types.ServiceDefinition) outputs.CompilerOutputLocation); ok {
		r0 = rf(serviceDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(outputs.CompilerOutputLocation)
		}
	}

	if rf, ok := ret.Get(1).(func(types.ServiceDefinition) error); ok {
		r1 = rf(serviceDef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateModelOutput provides a mock function with given fields: model
func (_m *MockCompilerOutputsManager) CreateModelOutput(model types.EntitySpec) (outputs.CompilerOutputWriter, error) {
	ret := _m.Called(model)

	if len(ret) == 0 {
		panic("no return value specified for CreateModelOutput")
	}

	var r0 outputs.CompilerOutputWriter
	var r1 error
	if rf, ok := ret.Get(0).(func(types.EntitySpec) (outputs.CompilerOutputWriter, error)); ok {
		return rf(model)
	}
	if rf, ok := ret.Get(0).(func(types.EntitySpec) outputs.CompilerOutputWriter); ok {
		r0 = rf(model)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(outputs.CompilerOutputWriter)
		}
	}

	if rf, ok := ret.Get(1).(func(types.EntitySpec) error); ok {
		r1 = rf(model)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateServiceOutput provides a mock function with given fields: serviceDef
func (_m *MockCompilerOutputsManager) CreateServiceOutput(serviceDef types.ServiceDefinition) (outputs.CompilerOutputWriter, error) {
	ret := _m.Called(serviceDef)

	if len(ret) == 0 {
		panic("no return value specified for CreateServiceOutput")
	}

	var r0 outputs.CompilerOutputWriter
	var r1 error
	if rf, ok := ret.Get(0).(func(types.ServiceDefinition) (outputs.CompilerOutputWriter, error)); ok {
		return rf(serviceDef)
	}
	if rf, ok := ret.Get(0).(func(types.ServiceDefinition) outputs.CompilerOutputWriter); ok {
		r0 = rf(serviceDef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(outputs.CompilerOutputWriter)
		}
	}

	if rf, ok := ret.Get(1).(func(types.ServiceDefinition) error); ok {
		r1 = rf(serviceDef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareOutputDirectory provides a mock function with given fields: path
func (_m *MockCompilerOutputsManager) PrepareOutputDirectory(path string) error {
	ret := _m.Called(path)

	if len(ret) == 0 {
		panic("no return value specified for PrepareOutputDirectory")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockCompilerOutputsManager creates a new instance of MockCompilerOutputsManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCompilerOutputsManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCompilerOutputsManager {
	mock := &MockCompilerOutputsManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
