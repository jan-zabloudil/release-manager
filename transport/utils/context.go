package utils

import (
	"context"
	"net/http"

	svcmodel "release-manager/service/model"
	"release-manager/transport/model"

	"github.com/google/uuid"
)

type contextKey string

const (
	userContextKey          = contextKey("user")
	projectContextKey       = contextKey("project")
	projectMemberContextKey = contextKey("project_member")
	appContextKey           = contextKey("app")
	releaseContextKey       = contextKey("release")
)

func ContextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
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
	user, ok := r.Context().Value(projectContextKey).(*model.Project)
	if !ok {
		panic("missing project value in request context")
	}
	return user
}

func ContextProjectID(r *http.Request) uuid.UUID {
	return ContextProject(r).ID
}

func ContextSetProjectMember(r *http.Request, member svcmodel.ProjectMember) *http.Request {
	ctx := context.WithValue(r.Context(), projectMemberContextKey, member)
	return r.WithContext(ctx)
}

func ContextProjectMember(r *http.Request) svcmodel.ProjectMember {
	role, ok := r.Context().Value(projectMemberContextKey).(svcmodel.ProjectMember)
	if !ok {
		panic("missing project member in request context")
	}

	return role
}

func ContextSetApp(r *http.Request, p *svcmodel.App) *http.Request {
	ctx := context.WithValue(r.Context(), appContextKey, p)
	return r.WithContext(ctx)
}

func ContextApp(r *http.Request) *svcmodel.App {
	app, ok := r.Context().Value(appContextKey).(*svcmodel.App)
	if !ok {
		panic("missing app in request context")
	}
	return app
}

func ContextProjectIDFromApp(r *http.Request) uuid.UUID {
	return ContextApp(r).ProjectID
}

func ContextSetRelease(r *http.Request, rls *svcmodel.Release) *http.Request {
	ctx := context.WithValue(r.Context(), releaseContextKey, rls)
	return r.WithContext(ctx)
}

func ContextRelease(r *http.Request) *svcmodel.Release {
	app, ok := r.Context().Value(releaseContextKey).(*svcmodel.Release)
	if !ok {
		panic("missing release value in request context")
	}
	return app
}
