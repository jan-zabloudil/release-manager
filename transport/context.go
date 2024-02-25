package transport

import (
	"context"
	"net/http"

	"release-manager/transport/model"
)

type contextKey string

const (
	userContextKey    = contextKey("user")
	projectContextKey = contextKey("project")
)

func ContextSetUser(r *http.Request, user *model.AuthUser) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextUser(r *http.Request) *model.AuthUser {
	user, ok := r.Context().Value(userContextKey).(*model.AuthUser)
	if !ok {
		panic("missing auth user value in request context")
	}
	return user
}

func ContextSetProject(r *http.Request, p *model.Project) *http.Request {
	ctx := context.WithValue(r.Context(), projectContextKey, p)
	return r.WithContext(ctx)
}

func ContextProject(r *http.Request) *model.Project {
	project, ok := r.Context().Value(projectContextKey).(*model.Project)
	if !ok {
		panic("missing project value in request context")
	}
	return project
}
