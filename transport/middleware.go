package transport

import (
	"errors"
	"net/http"
	"strings"

	reperr "release-manager/repository/errors"
	svcmodel "release-manager/service/model"
	neterr "release-manager/transport/errors"
	"release-manager/transport/model"

	"github.com/google/uuid"
	httpx "go.strv.io/net/http"
)

type ProjectIDFunc func(r *http.Request) uuid.UUID

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get(httpx.Header.Authorization)

		if authorizationHeader == "" {
			r = ContextSetUser(r, model.AnonUser)
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

		u := model.ToNetUser(
			user.ID,
			user.IsAdmin,
			user.Email,
			user.Name,
			user.AvatarUrl,
			user.CreatedAt,
			user.UpdatedAt,
		)
		r = ContextSetUser(r, &u)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextUser(r)

		if user.IsAnon() {
			WriteInvalidAuthenticationResponse(w, neterr.ErrAccessDeniedToAnonUser)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requireAdminUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextUser(r)

		if !user.IsAdmin {
			WriteForbiddenErrorResponse(w, neterr.ErrAccessDeniedToNonAdminUser)
			return
		}

		next.ServeHTTP(w, r)
	})

	return h.requireAuthUser(fn)
}

func (h *Handler) requireProjectMemberRole(f ProjectIDFunc, role svcmodel.ProjectRole, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := ContextUser(r)

		projectID := f(r)
		m, err := h.ProjectMemberSvc.Get(r.Context(), projectID, u.ID)
		if err != nil {
			switch {
			case errors.Is(err, reperr.ErrResourceNotFound):
				// User is not a member, do not expose resource and return 404 (instead of 403)
				WriteNotFoundResponse(w, err)
				return
			default:
				WriteServerErrorResponse(w, err)
				return
			}
		}

		if !m.HasAtLeastRole(role) {
			WriteForbiddenErrorResponse(w, neterr.ErrInsufficientProjectRole)
			return
		}

		r = ContextSetProjectMember(r, m)

		next.ServeHTTP(w, r)
	})

	return h.requireAuthUser(fn)
}

func (h *Handler) handleProject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := GetUUIDParamFrom(r, "id")
		if err != nil {
			WriteNotFoundResponse(w, err)
			return
		}

		p, err := h.ProjectSvc.Get(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, reperr.ErrResourceNotFound):
				WriteNotFoundResponse(w, err)
				return
			default:
				WriteServerErrorResponse(w, err)
				return
			}
		}
		netp := model.ToNetProject(p)
		r = ContextSetProject(r, &netp)

		next.ServeHTTP(w, r)
	})
}
