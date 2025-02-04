package domain

// EventFeedChanged представляет структуру события для Kafka.
type EventFeedChanged struct {
	UserID UserKey `json:"user_id"`
}
