package domain

import (
	"time"
)

type PostKey int64
type PostText string

// EventType определяет тип события.
type EventType string

// Возможные значения EventType.
const (
	EventPostCreated EventType = "post_created" // Пост создан
	EventPostEdited  EventType = "post_edited"  // Пост отредактирован
	EventPostDeleted EventType = "post_deleted" // Пост удален
)

type EventPostModified struct {
	Event           EventType
	Post            *Post     `json:"post"`
	AntiLadyGagaIds []UserKey `json:"anti_lady_gaga_ids"`
}

type Post struct {
	Id         PostKey    `json:"id"`
	UserId     UserKey    `json:"user_id"`
	Message    PostText   `json:"text"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at"`
}

type EventFollowerNotifyContent struct {
	Post  *Post     `json:"post"`
	Event EventType `json:"event"`
}

type EventFollowerNotify struct {
	Recipient UserKey                     `json:"recipient"`
	Content   *EventFollowerNotifyContent `json:"content"`
}
