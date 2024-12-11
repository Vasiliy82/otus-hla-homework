// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/Vasiliy82/otus-hla-homework/domain"
	jwt "github.com/golang-jwt/jwt/v5"

	mock "github.com/stretchr/testify/mock"
)

// JWTService is an autogenerated mock type for the JWTService type
type JWTService struct {
	mock.Mock
}

// ExtractClaims provides a mock function with given fields: token
func (_m *JWTService) ExtractClaims(token *jwt.Token) (*domain.UserClaims, error) {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for ExtractClaims")
	}

	var r0 *domain.UserClaims
	var r1 error
	if rf, ok := ret.Get(0).(func(*jwt.Token) (*domain.UserClaims, error)); ok {
		return rf(token)
	}
	if rf, ok := ret.Get(0).(func(*jwt.Token) *domain.UserClaims); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserClaims)
		}
	}

	if rf, ok := ret.Get(1).(func(*jwt.Token) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateToken provides a mock function with given fields: userID, permissions
func (_m *JWTService) GenerateToken(userID domain.UserKey, permissions []domain.Permission) (domain.TokenString, error) {
	ret := _m.Called(userID, permissions)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToken")
	}

	var r0 domain.TokenString
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserKey, []domain.Permission) (domain.TokenString, error)); ok {
		return rf(userID, permissions)
	}
	if rf, ok := ret.Get(0).(func(domain.UserKey, []domain.Permission) domain.TokenString); ok {
		r0 = rf(userID, permissions)
	} else {
		r0 = ret.Get(0).(domain.TokenString)
	}

	if rf, ok := ret.Get(1).(func(domain.UserKey, []domain.Permission) error); ok {
		r1 = rf(userID, permissions)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RevokeToken provides a mock function with given fields: token
func (_m *JWTService) RevokeToken(token *jwt.Token) error {
	ret := _m.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for RevokeToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*jwt.Token) error); ok {
		r0 = rf(token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateToken provides a mock function with given fields: tokenString
func (_m *JWTService) ValidateToken(tokenString domain.TokenString) (*jwt.Token, error) {
	ret := _m.Called(tokenString)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *jwt.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.TokenString) (*jwt.Token, error)); ok {
		return rf(tokenString)
	}
	if rf, ok := ret.Get(0).(func(domain.TokenString) *jwt.Token); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.TokenString) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJWTService creates a new instance of JWTService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJWTService(t interface {
	mock.TestingT
	Cleanup(func())
}) *JWTService {
	mock := &JWTService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}