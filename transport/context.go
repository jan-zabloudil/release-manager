package transport

import (
	"context"
	"net/http"

	"github.com/jan-zabloudil/release-manager/transport/model"
)

type contextKey string

const userContextKey = contextKey("user")

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
