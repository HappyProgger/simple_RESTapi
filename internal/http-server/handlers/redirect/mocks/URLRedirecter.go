// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// URLRedirecter is an autogenerated mock type for the URLRedirecter type
type URLRedirecter struct {
	mock.Mock
}

// GetURL provides a mock function with given fields: alias
func (_m *URLRedirecter) GetURL(alias string) (string, error) {
	ret := _m.Called(alias)

	if len(ret) == 0 {
		panic("no return value specified for GetURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewURLRedirecter creates a new instance of URLRedirecter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLRedirecter(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLRedirecter {
	mock := &URLRedirecter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
