package middleware

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func HandleProject(projectSvc model.ProjectService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := utils.GetUUIDParamFrom(r, "id")
			if err != nil {
				utils.WriteNotFoundResponse(w, err)
				return
			}

			p, err := projectSvc.Get(r.Context(), id)
			if err != nil {
				switch {
				case errors.Is(err, reperr.ErrResourceNotFound):
					utils.WriteNotFoundResponse(w, err)
					return
				default:
					utils.WriteServerErrorResponse(w, err)
					return
				}
			}
			netp := model.ToNetProject(p)
			r = utils.ContextSetProject(r, &netp)

			next.ServeHTTP(w, r)
		})
	}
}
