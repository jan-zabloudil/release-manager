package middleware

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	svcmodel "release-manager/service/model"
	neterr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"

	"github.com/google/uuid"
)

type ProjectIDFunc func(r *http.Request) uuid.UUID

func RequireAuthUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.ContextUser(r)

		if user.IsAnon() {
			utils.WriteInvalidAuthenticationResponse(w, neterr.ErrAccessDeniedToAnonUser)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireAdminUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.ContextUser(r)

		if !user.IsAdmin {
			utils.WriteForbiddenErrorResponse(w, neterr.ErrAccessDeniedToNonAdminUser)
			return
		}

		next.ServeHTTP(w, r)
	})

	return RequireAuthUser(fn)
}

func RequireProjectMemberRole(
	projectMemberSvc model.ProjectMemberService,
	f ProjectIDFunc,
	role svcmodel.ProjectRole,
	next http.HandlerFunc,
) http.HandlerFunc {

	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := utils.ContextUser(r)

		projectID := f(r)
		m, err := projectMemberSvc.Get(r.Context(), projectID, u.ID)
		if err != nil {
			switch {
			case errors.Is(err, reperr.ErrResourceNotFound):
				// User is not a member, do not expose resource and return 404 (instead of 403)
				utils.WriteNotFoundResponse(w, err)
				return
			default:
				utils.WriteServerErrorResponse(w, err)
				return
			}
		}

		if !m.HasAtLeastRole(role) {
			utils.WriteForbiddenErrorResponse(w, neterr.ErrInsufficientProjectRole)
			return
		}

		r = utils.ContextSetProjectMember(r, m)

		next.ServeHTTP(w, r)
	})

	return RequireAuthUser(fn)
}
