package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	authUserIDContextKey    contextKey = "auth_user_id"
	userIDContextKey        contextKey = "user_id"
	projectIDContextKey     contextKey = "project_id"
	environmentIDContextKey contextKey = "environment_id"
)

type ContextSetUUIDFunc func(r *http.Request, id uuid.UUID) *http.Request

func ContextProjectID(r *http.Request) uuid.UUID {
	return contextUUID(r, projectIDContextKey)
}

func ContextSetProjectID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, projectIDContextKey)
}

func ContextAuthUserID(r *http.Request) uuid.UUID {
	return contextUUID(r, authUserIDContextKey)
}

func ContextSetAuthUserID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, authUserIDContextKey)
}

func ContextUserID(r *http.Request) uuid.UUID {
	return contextUUID(r, userIDContextKey)
}

func ContextSetUserID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, userIDContextKey)
}

func ContextEnvironmentID(r *http.Request) uuid.UUID {
	return contextUUID(r, environmentIDContextKey)
}

func ContextSetEnvironmentID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, environmentIDContextKey)
}

func contextSetUUID(r *http.Request, id uuid.UUID, key contextKey) *http.Request {
	ctx := context.WithValue(r.Context(), key, id)
	return r.WithContext(ctx)
}

func contextUUID(r *http.Request, key contextKey) uuid.UUID {
	user, ok := r.Context().Value(key).(uuid.UUID)
	if !ok {
		panic(fmt.Sprintf("missing %s value in request context", key))
	}
	return user
}
