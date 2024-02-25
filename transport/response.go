package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Response map[string]any

func logErrResponse(r *http.Request, err error, level slog.Level) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	slog.Log(context.TODO(), level, err.Error(), "method", method, "uri", uri)
}

func WriteResponse(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func WriteNotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	logErrResponse(r, err, slog.LevelDebug)

	message := "the requested resource could not be found"
	writeErrorResponse(w, r, http.StatusNotFound, message)
}

func WriteMethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	logErrResponse(r, err, slog.LevelDebug)

	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	writeErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func WriteServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logErrResponse(r, err, slog.LevelError)

	message := "the server encountered a problem and could not process your request."
	writeErrorResponse(w, r, http.StatusInternalServerError, message)
}

func WriteInvalidAuthenticationResponse(w http.ResponseWriter, r *http.Request, err error) {
	logErrResponse(r, err, slog.LevelDebug)

	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	writeErrorResponse(w, r, http.StatusUnauthorized, message)
}

func WriteForbiddenErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logErrResponse(r, err, slog.LevelDebug)

	message := "you do not have permission to perform this action"
	writeErrorResponse(w, r, http.StatusForbidden, message)
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Response{"error": message}

	err := WriteResponse(w, status, env, nil)
	if err != nil {
		logErrResponse(r, err, slog.LevelError)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
