package transport

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type Response map[string]any

func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	slog.Error(err.Error(), "method", method, "uri", uri)
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

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	if err := WriteResponse(w, status, Response{"error": message}, nil); err != nil {
		logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	WriteErrorResponse(w, r, http.StatusNotFound, message)
}

func WriteMethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	WriteErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func WriteServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "The server encountered a problem and could not process your request."
	WriteErrorResponse(w, r, http.StatusInternalServerError, message)
}
