// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// PostRepository is an autogenerated mock type for the PostRepository type
type PostRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: userId, message
func (_m *PostRepository) Create(userId domain.UserKey, message domain.PostMessage) (domain.PostKey, error) {
	ret := _m.Called(userId, message)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 domain.PostKey
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserKey, domain.PostMessage) (domain.PostKey, error)); ok {
		return rf(userId, message)
	}
	if rf, ok := ret.Get(0).(func(domain.UserKey, domain.PostMessage) domain.PostKey); ok {
		r0 = rf(userId, message)
	} else {
		r0 = ret.Get(0).(domain.PostKey)
	}

	if rf, ok := ret.Get(1).(func(domain.UserKey, domain.PostMessage) error); ok {
		r1 = rf(userId, message)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *PostRepository) Delete(id domain.PostKey) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.PostKey) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: postId
func (_m *PostRepository) Get(postId domain.PostKey) (*domain.Post, error) {
	ret := _m.Called(postId)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.PostKey) (*domain.Post, error)); ok {
		return rf(postId)
	}
	if rf, ok := ret.Get(0).(func(domain.PostKey) *domain.Post); ok {
		r0 = rf(postId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.PostKey) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFeed provides a mock function with given fields: userId, limit
func (_m *PostRepository) GetFeed(userId domain.UserKey, limit int) ([]*domain.Post, error) {
	ret := _m.Called(userId, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetFeed")
	}

	var r0 []*domain.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserKey, int) ([]*domain.Post, error)); ok {
		return rf(userId, limit)
	}
	if rf, ok := ret.Get(0).(func(domain.UserKey, int) []*domain.Post); ok {
		r0 = rf(userId, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.UserKey, int) error); ok {
		r1 = rf(userId, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPostOwner provides a mock function with given fields: postId
func (_m *PostRepository) GetPostOwner(postId domain.PostKey) (domain.UserKey, error) {
	ret := _m.Called(postId)

	if len(ret) == 0 {
		panic("no return value specified for GetPostOwner")
	}

	var r0 domain.UserKey
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.PostKey) (domain.UserKey, error)); ok {
		return rf(postId)
	}
	if rf, ok := ret.Get(0).(func(domain.PostKey) domain.UserKey); ok {
		r0 = rf(postId)
	} else {
		r0 = ret.Get(0).(domain.UserKey)
	}

	if rf, ok := ret.Get(1).(func(domain.PostKey) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: userId, limit, lastPostId
func (_m *PostRepository) List(userId domain.UserKey, limit int, lastPostId domain.PostKey) ([]*domain.Post, error) {
	ret := _m.Called(userId, limit, lastPostId)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*domain.Post
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.UserKey, int, domain.PostKey) ([]*domain.Post, error)); ok {
		return rf(userId, limit, lastPostId)
	}
	if rf, ok := ret.Get(0).(func(domain.UserKey, int, domain.PostKey) []*domain.Post); ok {
		r0 = rf(userId, limit, lastPostId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Post)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.UserKey, int, domain.PostKey) error); ok {
		r1 = rf(userId, limit, lastPostId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMessage provides a mock function with given fields: postId, newMessage
func (_m *PostRepository) UpdateMessage(postId domain.PostKey, newMessage domain.PostMessage) error {
	ret := _m.Called(postId, newMessage)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(domain.PostKey, domain.PostMessage) error); ok {
		r0 = rf(postId, newMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPostRepository creates a new instance of PostRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPostRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PostRepository {
	mock := &PostRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
