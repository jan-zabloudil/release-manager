package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"release-manager/transport/utils"

	"github.com/go-playground/validator/v10"
	httpx "go.strv.io/net/http"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	if err := httpx.WriteResponse(w, data, status, httpx.WithContentType(httpx.ApplicationJSON)); err != nil {
		slog.Error("writing json response", "error", err)
	}
}

func WriteErrorResponse(w http.ResponseWriter, status int, opts ...httpx.ErrorResponseOption) {
	if err := httpx.WriteErrorResponse(w, status, opts...); err != nil {
		slog.Error("writing error response", "error", err)
	}
}

func WriteNotFoundResponse(w http.ResponseWriter, err error) {
	msg := "the requested resource could not be found"
	WriteErrorResponse(
		w,
		http.StatusNotFound,
		httpx.WithError(err),
		httpx.WithErrorCode("404"),
		httpx.WithErrorMessage(msg),
	)
}

func WriteMethodNotAllowedResponse(w http.ResponseWriter, err error) {
	msg := "the method is not supported for this resource"
	WriteErrorResponse(
		w,
		http.StatusMethodNotAllowed,
		httpx.WithError(err),
		httpx.WithErrorCode(strconv.Itoa(http.StatusMethodNotAllowed)),
		httpx.WithErrorMessage(msg),
	)
}

func WriteServerErrorResponse(w http.ResponseWriter, err error) {
	msg := "the server encountered a problem and could not process your request."
	WriteErrorResponse(
		w,
		http.StatusInternalServerError,
		httpx.WithError(err),
		httpx.WithErrorCode(strconv.Itoa(http.StatusInternalServerError)),
		httpx.WithErrorMessage(msg),
	)
}

func WriteInvalidAuthenticationResponse(w http.ResponseWriter, err error) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	msg := "invalid or missing authentication token"
	WriteErrorResponse(
		w,
		http.StatusUnauthorized,
		httpx.WithError(err),
		httpx.WithErrorCode(strconv.Itoa(http.StatusUnauthorized)),
		httpx.WithErrorMessage(msg),
	)
}

func WriteForbiddenErrorResponse(w http.ResponseWriter, err error) {
	msg := "you do not have permission to perform this action"
	WriteErrorResponse(
		w,
		http.StatusForbidden,
		httpx.WithError(err),
		httpx.WithErrorCode(strconv.Itoa(http.StatusForbidden)),
		httpx.WithErrorMessage(msg),
	)
}

func WriteBadRequestResponse(w http.ResponseWriter, err error) {
	WriteErrorResponse(
		w,
		http.StatusBadRequest,
		httpx.WithError(err),
		httpx.WithErrorCode(strconv.Itoa(http.StatusBadRequest)),
		httpx.WithErrorMessage(err.Error()),
	)
}

func WriteUnprocessableEntityResponse(w http.ResponseWriter, err error) {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		msg := "validation errors"
		te := utils.TranslateValidationErrs(verrs)
		WriteErrorResponse(
			w,
			http.StatusUnprocessableEntity,
			httpx.WithError(err),
			httpx.WithErrorMessage(msg),
			httpx.WithErrorCode(strconv.Itoa(http.StatusUnprocessableEntity)),
			httpx.WithErrorData(te),
		)

		return
	}

	WriteErrorResponse(
		w,
		http.StatusUnprocessableEntity,
		httpx.WithError(err),
		httpx.WithErrorCode(strconv.Itoa(http.StatusUnprocessableEntity)),
		httpx.WithErrorMessage(err.Error()),
	)
}
