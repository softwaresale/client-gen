// Code generated by mockery v2.46.3. DO NOT EDIT.

package outputs_mocks

import mock "github.com/stretchr/testify/mock"

// MockCompilerOutputWriter is an autogenerated mock type for the CompilerOutputWriter type
type MockCompilerOutputWriter struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockCompilerOutputWriter) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Location provides a mock function with given fields:
func (_m *MockCompilerOutputWriter) Location() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Location")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *MockCompilerOutputWriter) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Write provides a mock function with given fields: p
func (_m *MockCompilerOutputWriter) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockCompilerOutputWriter creates a new instance of MockCompilerOutputWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCompilerOutputWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCompilerOutputWriter {
	mock := &MockCompilerOutputWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
