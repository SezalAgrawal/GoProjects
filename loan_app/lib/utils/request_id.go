package utils

import (
	"context"
)

type logKey string

const (
	RequestIDHeader              = "X-Request-Id"
	RequestIDLogKey       logKey = "request_id"
	CurrentUserIDStoreKey        = "X_CURRENT_USER_ID"
)

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDLogKey, reqID)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDLogKey).(string); ok {
		return reqID
	}
	return ""
}
