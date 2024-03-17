package middleware

import (
	"errors"
	"net/http"
	"strings"

	reperr "release-manager/repository/errors"
	neterr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"

	httpx "go.strv.io/net/http"
)

func Auth(userSvc model.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")
			authorizationHeader := r.Header.Get(httpx.Header.Authorization)

			if authorizationHeader == "" {
				r = utils.ContextSetUser(r, model.AnonUser)
				next.ServeHTTP(w, r)
				return
			}

			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				utils.WriteInvalidAuthenticationResponse(w, neterr.ErrInvalidBearer)
				return
			}

			tokenString := headerParts[1]

			user, err := userSvc.GetForToken(r.Context(), tokenString)
			if err != nil {
				switch {
				case errors.Is(err, reperr.ErrUserAuthenticationFailed):
					utils.WriteInvalidAuthenticationResponse(w, err)
					return
				default:
					utils.WriteServerErrorResponse(w, err)
					return
				}
			}

			u := model.ToNetUser(
				user.ID,
				user.IsAdmin,
				user.Email,
				user.Name,
				user.AvatarUrl,
				user.CreatedAt,
				user.UpdatedAt,
			)
			r = utils.ContextSetUser(r, &u)

			next.ServeHTTP(w, r)
		})
	}
}
