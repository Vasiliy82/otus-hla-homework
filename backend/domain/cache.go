package domain

// EventInvalidateCache представляет структуру события для Kafka.
type EventInvalidateCache struct {
	UserID    UserKey   `json:"user_id"`
	EventType EventType `json:"event_type"`
}

// EventType определяет тип события.
type EventType string

// Возможные значения EventType.
const (
	EventPostCreated EventType = "post_created" // Пост создан
	EventPostEdited  EventType = "post_edited"  // Пост отредактирован
	EventPostDeleted EventType = "post_deleted" // Пост удален
	EventFeedRefresh EventType = "feed_refresh" // Лента нуждается в пересчете
)
