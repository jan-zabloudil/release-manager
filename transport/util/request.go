package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpx "go.strv.io/net/http"
	"go.strv.io/net/http/param"
)

const (
	SignatureHeader = "X-Hub-Signature-256"
	// GithubHookEvent is the header key for the GitHub webhook event type.
	// Docs: https://docs.github.com/en/webhooks/webhook-events-and-payloads
	GithubHookEvent = "X-GitHub-Event"
)

// UnmarshalBody unmarshals the request body into the provided struct.
func UnmarshalBody(r *http.Request, b any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, b); err != nil {
		return err
	}

	return nil
}

func RequestID(h http.Header) string {
	return h.Get(httpx.Header.XRequestID)
}

func GetUUIDFromURL(r *http.Request, key string) (uuid.UUID, error) {
	idFromURL := chi.URLParam(r, key)

	id, err := uuid.Parse(idFromURL)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

// UnmarshalURLParams parses URL path and query parameters into a struct with tagged fields.
func UnmarshalURLParams[TParams any](r *http.Request) (TParams, error) {
	var params TParams
	if err := param.DefaultParser().Parse(r, &params); err != nil {
		return params, err
	}
	return params, nil
}

type ParamUnmarshaller interface {
	UnmarshalText(data []byte) error
}

// GetQueryParam retrieves a URL query parameter from the request and unmarshals it into the provided type.
func GetQueryParam[TParam any, TPtrParam interface {
	*TParam
	ParamUnmarshaller
}](r *http.Request, paramName string) (TParam, error) {
	paramValue := r.URL.Query().Get(paramName)
	return unmarshalParam[TParam, TPtrParam](paramValue)
}

// GetPathParam retrieves a URL path parameter from the request and unmarshals it into the provided type.
func GetPathParam[TParam any, TPtrParam interface {
	*TParam
	ParamUnmarshaller
}](r *http.Request, paramName string) (TParam, error) {
	paramValue := chi.URLParam(r, paramName)
	return unmarshalParam[TParam, TPtrParam](paramValue)
}

func unmarshalParam[TParam any, TPtrParam interface {
	*TParam
	ParamUnmarshaller
}](paramValue string) (TParam, error) {
	var zeroValue TParam
	p := TPtrParam(new(TParam))
	if err := p.UnmarshalText([]byte(paramValue)); err != nil {
		return zeroValue, err
	}
	return *p, nil
}
