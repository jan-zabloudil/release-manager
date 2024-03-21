package transport

import (
	"errors"
	"net/http"
	"strings"

	reperr "release-manager/repository/errors"
	svcmodel "release-manager/service/model"
	neterr "release-manager/transport/errors"
	"release-manager/transport/utils"

	httpx "go.strv.io/net/http"
)

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header.Get(httpx.Header.Authorization)

		if authorizationHeader == "" {
			r = utils.ContextSetUser(r, svcmodel.AnonUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			WriteInvalidAuthenticationResponse(w, neterr.ErrInvalidBearer)
			return
		}

		tokenString := headerParts[1]

		user, err := h.UserSvc.GetForToken(r.Context(), tokenString)
		if err != nil {
			switch {
			case errors.Is(err, reperr.ErrUserAuthenticationFailed):
				WriteInvalidAuthenticationResponse(w, err)
				return
			default:
				WriteServerErrorResponse(w, err)
				return
			}
		}

		r = utils.ContextSetUser(r, &user)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.ContextUser(r)

		if user.IsAnon() {
			WriteInvalidAuthenticationResponse(w, neterr.ErrAccessDeniedToAnonUser)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAdminUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.ContextUser(r)

		if !user.IsAdmin() {
			WriteForbiddenErrorResponse(w, neterr.ErrAccessDeniedToNonAdminUser)
			return
		}

		next.ServeHTTP(w, r)
	})

	return h.requireAuthUser(fn)
}
