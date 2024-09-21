package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"release-manager/auth"
	cryptox "release-manager/pkg/crypto"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/util"

	"github.com/google/uuid"
	httpx "go.strv.io/net/http"
)

type AuthClient interface {
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}

func Auth(authClient AuthClient) func(next http.Handler) http.Handler {
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
				util.WriteResponseError(w, resperrors.NewNotBearerTokenFormatError())
				return
			}

			tokenString := headerParts[1]

			id, err := authClient.Authenticate(r.Context(), tokenString)
			if err != nil {
				if errors.Is(err, auth.ErrInvalidOrExpiredToken) {
					util.WriteResponseError(w, resperrors.NewExpiredOrInvalidTokenError().Wrap(err))
					return
				}

				util.WriteResponseError(w, resperrors.NewServerError().Wrap(err))
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
			util.WriteResponseError(w, resperrors.NewMissingBearerTokenError())
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SetResourceUUIDToContext is a middleware that extracts resource ID (UUID) from the URL and sets it to the request context.
func SetResourceUUIDToContext(idKey string, f util.ContextSetUUIDFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := util.GetUUIDFromURL(r, idKey)
			if err != nil {
				util.WriteResponseError(w, resperrors.NewInvalidResourceIDError().Wrap(err))
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
