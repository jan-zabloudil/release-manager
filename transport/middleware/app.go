package middleware

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func HandleApp(appSvc model.AppService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id, err := utils.GetUUIDParamFrom(r, "appId")
			if err != nil {
				utils.WriteNotFoundResponse(w, err)
				return
			}

			app, err := appSvc.Get(r.Context(), id)
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

			r = utils.ContextSetApp(r, &app)

			next.ServeHTTP(w, r)
		})
	}
}
