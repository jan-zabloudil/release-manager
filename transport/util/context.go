package util

import (
	"context"
	"fmt"
	"net/http"

	"release-manager/pkg/id"

	"github.com/google/uuid"
)

type contextKey string

const (
	authUserIDContextKey contextKey = "auth_user_id"
	projectIDContextKey  contextKey = "project_id"
)

type ContextSetUUIDFunc func(r *http.Request, id uuid.UUID) *http.Request

func ContextProjectID(r *http.Request) uuid.UUID {
	return contextUUID(r, projectIDContextKey)
}

func ContextSetProjectID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, projectIDContextKey)
}

func contextSetUUID(r *http.Request, id uuid.UUID, key contextKey) *http.Request {
	return setContextValue(r, key, id)
}

func contextUUID(r *http.Request, key contextKey) uuid.UUID {
	user, ok := r.Context().Value(key).(uuid.UUID)
	if !ok {
		panic(fmt.Sprintf("missing %s value in request context", key))
	}
	return user
}

// Refactoring in progress
// Methods for context values with custom type instead of general UUID

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
