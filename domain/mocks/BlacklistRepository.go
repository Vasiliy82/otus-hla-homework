// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// BlacklistRepository is an autogenerated mock type for the BlacklistRepository type
type BlacklistRepository struct {
	mock.Mock
}

// AddToBlacklist provides a mock function with given fields: serial, expireDate
func (_m *BlacklistRepository) AddToBlacklist(serial int64, expireDate time.Time) error {
	ret := _m.Called(serial, expireDate)

	if len(ret) == 0 {
		panic("no return value specified for AddToBlacklist")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, time.Time) error); ok {
		r0 = rf(serial, expireDate)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsBlacklisted provides a mock function with given fields: serial
func (_m *BlacklistRepository) IsBlacklisted(serial int64) (bool, error) {
	ret := _m.Called(serial)

	if len(ret) == 0 {
		panic("no return value specified for IsBlacklisted")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (bool, error)); ok {
		return rf(serial)
	}
	if rf, ok := ret.Get(0).(func(int64) bool); ok {
		r0 = rf(serial)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(serial)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSerial provides a mock function with given fields:
func (_m *BlacklistRepository) NewSerial() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewSerial")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBlacklistRepository creates a new instance of BlacklistRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBlacklistRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BlacklistRepository {
	mock := &BlacklistRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
