package util

import (
	"log/slog"
	"net/http"

	resperrors "release-manager/transport/errors"

	httpx "go.strv.io/net/http"
)

func WriteResponseError(w http.ResponseWriter, r *resperrors.Error) {
	if err := httpx.WriteErrorResponse(
		w,
		r.StatusCode,
		httpx.WithError(r.Err),
		httpx.WithErrorCode(r.Code),
		httpx.WithErrorMessage(r.Message),
	); err != nil {
		slog.Error("writing error response", "error", err)
	}
}

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	if err := httpx.WriteResponse(w, data, status, httpx.WithContentType(httpx.ApplicationJSON)); err != nil {
		slog.Error("writing json response", "error", err)
	}
}
