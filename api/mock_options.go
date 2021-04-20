// Code generated by mockery v2.6.0. DO NOT EDIT.

package api

import mock "github.com/stretchr/testify/mock"

// MockOptions is an autogenerated mock type for the Options type
type MockOptions struct {
	mock.Mock
}

// AllKeys provides a mock function with given fields:
func (_m *MockOptions) AllKeys() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Get provides a mock function with given fields: _a0
func (_m *MockOptions) Get(_a0 string) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// GetBool provides a mock function with given fields: _a0
func (_m *MockOptions) GetBool(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetInt provides a mock function with given fields: _a0
func (_m *MockOptions) GetInt(_a0 string) int {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetString provides a mock function with given fields: _a0
func (_m *MockOptions) GetString(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Set provides a mock function with given fields: _a0, _a1
func (_m *MockOptions) Set(_a0 string, _a1 interface{}) {
	_m.Called(_a0, _a1)
}