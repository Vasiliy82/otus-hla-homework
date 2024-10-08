// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/Vasiliy82/otus-hla-homework/domain"
	jwt "github.com/golang-jwt/jwt/v5"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// AddFriend provides a mock function with given fields: my_id, friend_id
func (_m *UserService) AddFriend(my_id string, friend_id string) error {
	ret := _m.Called(my_id, friend_id)

	if len(ret) == 0 {
		panic("no return value specified for AddFriend")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(my_id, friend_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetById provides a mock function with given fields: id
func (_m *UserService) GetById(id string) (*domain.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: username, password
func (_m *UserService) Login(username string, password string) (domain.TokenString, error) {
	ret := _m.Called(username, password)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 domain.TokenString
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (domain.TokenString, error)); ok {
		return rf(username, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) domain.TokenString); ok {
		r0 = rf(username, password)
	} else {
		r0 = ret.Get(0).(domain.TokenString)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logout provides a mock function with given fields: token
func (_m *UserService) Logout(token *jwt.Token) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*jwt.Token) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterUser provides a mock function with given fields: user
func (_m *UserService) RegisterUser(user *domain.User) (string, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for RegisterUser")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*domain.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(*domain.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*domain.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFriend provides a mock function with given fields: my_id, friend_id
func (_m *UserService) RemoveFriend(my_id string, friend_id string) error {
	ret := _m.Called(my_id, friend_id)

	if len(ret) == 0 {
		panic("no return value specified for RemoveFriend")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(my_id, friend_id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: firstName, lastName
func (_m *UserService) Search(firstName string, lastName string) ([]*domain.User, error) {
	ret := _m.Called(firstName, lastName)

	if len(ret) == 0 {
		panic("no return value specified for Search")
	}

	var r0 []*domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]*domain.User, error)); ok {
		return rf(firstName, lastName)
	}
	if rf, ok := ret.Get(0).(func(string, string) []*domain.User); ok {
		r0 = rf(firstName, lastName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(firstName, lastName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
