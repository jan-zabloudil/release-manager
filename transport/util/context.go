package util

import (
	"context"
	"fmt"
	"net/http"

	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/id"

	"github.com/google/uuid"
)

type contextKey string

const (
	authUserIDContextKey             contextKey = "auth_user_id"
	projectIDContextKey              contextKey = "project_id"
	environmentIDContextKey          contextKey = "environment_id"
	projectInvitationIDContextKey    contextKey = "project_invitation_id"
	projectInvitationTokenContextKey contextKey = "project_invitation_token"
	releaseIDContextKey              contextKey = "release_id"
)

type ContextSetUUIDFunc func(r *http.Request, id uuid.UUID) *http.Request

func ContextProjectID(r *http.Request) uuid.UUID {
	return contextUUID(r, projectIDContextKey)
}

func ContextSetProjectID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, projectIDContextKey)
}

func ContextEnvironmentID(r *http.Request) uuid.UUID {
	return contextUUID(r, environmentIDContextKey)
}

func ContextSetEnvironmentID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, environmentIDContextKey)
}

func ContextProjectInvitationID(r *http.Request) uuid.UUID {
	return contextUUID(r, projectInvitationIDContextKey)
}

func ContextSetProjectInvitationID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, projectInvitationIDContextKey)
}

func ContextReleaseID(r *http.Request) uuid.UUID {
	return contextUUID(r, releaseIDContextKey)
}

func ContextSetReleaseID(r *http.Request, id uuid.UUID) *http.Request {
	return contextSetUUID(r, id, releaseIDContextKey)
}

func ContextSetProjectInvitationToken(r *http.Request, tkn cryptox.Token) *http.Request {
	ctx := context.WithValue(r.Context(), projectInvitationTokenContextKey, tkn)
	return r.WithContext(ctx)
}

func ContextProjectInvitationToken(r *http.Request) cryptox.Token {
	tkn, ok := r.Context().Value(projectInvitationTokenContextKey).(cryptox.Token)
	if !ok {
		panic(fmt.Sprintf("missing %s value in request context", projectInvitationTokenContextKey))
	}
	return tkn
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
