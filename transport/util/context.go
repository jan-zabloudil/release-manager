package util

import (
	"context"
	"fmt"
	"net/http"

	"release-manager/pkg/id"
)

type contextKey string

const (
	authUserIDContextKey contextKey = "auth_user_id"
)

func ContextAuthUserID(r *http.Request) id.AuthUser {
	return mustContextValue[id.AuthUser](r, authUserIDContextKey)
}

func ContextSetAuthUserID(r *http.Request, id id.AuthUser) *http.Request {
	return setContextValue(r, authUserIDContextKey, id)
}

func setContextValue(r *http.Request, key contextKey, value any) *http.Request {
	ctx := context.WithValue(r.Context(), key, value)
	return r.WithContext(ctx)
}

func mustContextValue[T any](r *http.Request, key contextKey) T {
	value, ok := r.Context().Value(key).(T)
	if !ok {
		panic(fmt.Sprintf("missing %s value in request context", key))
	}
	return value
}
