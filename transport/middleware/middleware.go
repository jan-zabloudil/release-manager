package middleware

import (
	"net/http"
	"strings"

	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"

	"github.com/google/uuid"
	httpx "go.strv.io/net/http"
)

func Auth(authSvc model.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(httpx.Header.Authorization)

			if authorizationHeader == "" {
				r = util.ContextSetAuthUserID(r, uuid.Nil)
				next.ServeHTTP(w, r)
				return
			}

			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				util.WriteResponseError(w, responseerrors.NewNotBearerTokenFormatError())
				return
			}

			tokenString := headerParts[1]

			id, err := authSvc.Authenticate(r.Context(), tokenString)
			if err != nil {
				util.WriteResponseError(w, responseerrors.NewUnauthorizedError().Wrap(err))
				return
			}

			r = util.ContextSetAuthUserID(r, id)

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := util.ContextAuthUserID(r); id == uuid.Nil {
			util.WriteResponseError(w, responseerrors.NewMissingBearerTokenError())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HandleResourceID(idKey string, f util.ContextSetUUIDFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := util.GetUUIDFromURL(r, idKey)
			if err != nil {
				util.WriteResponseError(w, responseerrors.NewInvalidResourceIDError().Wrap(err))
				return
			}

			r = f(r, id)

			next.ServeHTTP(w, r)
		})
	}
}

func HandleInvitationToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = util.ContextSetProjectInvitationToken(r, cryptox.Token(util.GetQueryParam(r, "token")))

		next.ServeHTTP(w, r)
	})
}
