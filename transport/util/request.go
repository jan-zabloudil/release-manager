package util

import (
	"encoding/json"
	"io"
	"net/http"

	"release-manager/pkg/validatorx"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
	"go.strv.io/net/http/param"
)

func RequestID(h http.Header) string {
	return h.Get(httpx.Header.XRequestID)
}

// UnmarshalBody unmarshals the request body into the provided struct and validates it.
func UnmarshalBody(r *http.Request, b any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, b); err != nil {
		return err
	}

	if err := validatorx.ValidateStruct(b); err != nil {
		return err
	}

	return nil
}

// UnmarshalURLParams parses URL path and query parameters into a struct and validates it.
func UnmarshalURLParams[TParams any](r *http.Request) (TParams, error) {
	var params TParams
	if err := param.DefaultParser().WithPathParamFunc(chi.URLParam).Parse(r, &params); err != nil {
		return params, err
	}

	if err := validatorx.ValidateStruct(params); err != nil {
		return params, err
	}

	return params, nil
}

type ParamUnmarshaller interface {
	UnmarshalText(data []byte) error
}

// GetPathParam retrieves a URL path parameter from the request and unmarshals it into the provided type.
func GetPathParam[TParam any, TPtrParam interface {
	*TParam
	ParamUnmarshaller
}](r *http.Request, paramName string) (pathParam TParam, err error) {
	p := TPtrParam(new(TParam))
	if err = p.UnmarshalText([]byte(chi.URLParam(r, paramName))); err != nil {
		return pathParam, err
	}
	return *p, nil
}
