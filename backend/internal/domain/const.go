package domain

// ContextKey для хранения x-request-id в контексте
type contextKey string

const RequestIDKey contextKey = "x-request-id"
const RequestIDHeader = "X-Request-ID"
