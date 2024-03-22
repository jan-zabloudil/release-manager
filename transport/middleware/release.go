package middleware

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func HandleRelease(rlsSvc model.ReleaseService, appSvc model.AppService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id, err := utils.GetUUIDParamFrom(r, "id")
			if err != nil {
				utils.WriteNotFoundResponse(w, err)
				return
			}

			rls, err := rlsSvc.Get(r.Context(), id)
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

			app, err := appSvc.Get(r.Context(), rls.AppID)
			if err != nil {
				utils.WriteServerErrorResponse(w, err)
				return
			}

			r = utils.ContextSetRelease(r, &rls)
			// Need to save app to context as well in order to be able to get project ID for authorization logic
			r = utils.ContextSetApp(r, &app)

			next.ServeHTTP(w, r)
		})
	}
}
