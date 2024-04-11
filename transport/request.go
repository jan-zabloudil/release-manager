package transport

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpx "go.strv.io/net/http"
)

func GetUUIDFromURL(r *http.Request, key string) (uuid.UUID, error) {
	idFromURL := chi.URLParam(r, key)

	id, err := uuid.Parse(idFromURL)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func UnmarshalRequest(r *http.Request, b any) error {
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
