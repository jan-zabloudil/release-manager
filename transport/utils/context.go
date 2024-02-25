package utils

import (
	"context"
	"net/http"

	svcmodel "release-manager/service/model"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *svcmodel.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextUser(r *http.Request) *svcmodel.User {
	user, ok := r.Context().Value(userContextKey).(*svcmodel.User)
	if !ok {
		panic("missing auth user value in request context")
	}
	return user
}
