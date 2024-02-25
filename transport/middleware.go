package transport

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	svcerr "github.com/jan-zabloudil/release-manager/service/errors"
	neterr "github.com/jan-zabloudil/release-manager/transport/errors"
	"github.com/jan-zabloudil/release-manager/transport/model"
)

func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				WriteServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = ContextSetUser(r, model.AnonUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			WriteInvalidAuthenticationResponse(w, r, neterr.ErrInvalidBearer)
			return
		}

		tokenString := headerParts[1]

		user, err := h.UserSvc.GetForToken(r.Context(), tokenString)
		if err != nil {
			switch {
			case errors.Is(err, svcerr.ErrUserAuthenticationFailed):
				WriteInvalidAuthenticationResponse(w, r, err)
				return
			default:
				WriteServerErrorResponse(w, r, err)
				return
			}
		}

		authUser := model.ToAuthUser(user)
		r = ContextSetUser(r, &authUser)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextUser(r)

		if user.IsAnon() {
			WriteInvalidAuthenticationResponse(w, r, neterr.ErrAccessDeniedToAnonUser)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAdminUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextUser(r)

		if !user.IsAdmin {
			WriteForbiddenErrorResponse(w, r, neterr.ErrAccessDeniedToNonAdminUser)
			return
		}

		next.ServeHTTP(w, r)
	})

	return h.requireAuthUser(fn)
}
