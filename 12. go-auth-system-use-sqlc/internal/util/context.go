package util

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	RequestIDKey contextKey = "request_id"
	StartTimeKey contextKey = "start_time"
)
