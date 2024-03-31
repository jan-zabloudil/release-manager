package transport

import (
	"net/http"
	"strings"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/util"

	"github.com/google/uuid"
	httpx "go.strv.io/net/http"
)

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(httpx.Header.Authorization)

		if authorizationHeader == "" {
			r = util.ContextSetAuthUserID(r, uuid.Nil)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			WriteResponseError(w, responseerrors.NewNotBearerTokenFormatError())
			return
		}

		tokenString := headerParts[1]

		id, err := h.AuthSvc.Authenticate(r.Context(), tokenString)
		if err != nil {
			WriteResponseError(w, responseerrors.NewUnauthorizedError().Wrap(err))
			return
		}

		r = util.ContextSetAuthUserID(r, id)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := util.ContextAuthUserID(r); id == uuid.Nil {
			WriteResponseError(w, responseerrors.NewMissingBearerTokenError())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) handleResourceID(idKey string, f util.ContextSetUUIDFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := GetUUIDFromURL(r, idKey)
			if err != nil {
				WriteResponseError(w, responseerrors.NewInvalidResourceIDError().Wrap(err))
				return
			}

			r = f(r, id)

			next.ServeHTTP(w, r)
		})
	}
}
