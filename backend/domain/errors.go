package domain

import "errors"

var (
	ErrObjectAlreadyExists = errors.New("object already exists")
	ErrObjectNotFound      = errors.New("object not found")

	ErrUserNotFound        = errors.New("user not found")
	ErrFriendNotFound      = errors.New("friend not found")
	ErrFriendAlreadyExists = errors.New("friend already exists")
)
