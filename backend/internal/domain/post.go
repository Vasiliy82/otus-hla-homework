package domain

import "time"

type PostKey int64
type PostMessage string

type Post struct {
	Id         PostKey
	UserId     UserKey
	Message    PostMessage
	CreatedAt  time.Time
	ModifiedAt *time.Time
}
