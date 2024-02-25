package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	Mux *chi.Mux
}

func NewHandler() *Handler {
	h := Handler{
		Mux: chi.NewRouter(),
	}
	h.Mux.Use(middleware.Logger) // TODO use slog
	h.Mux.Use(h.recoverPanic)

	h.Mux.Get("/ping", h.ping)

	h.Mux.NotFound(WriteNotFoundResponse)
	h.Mux.MethodNotAllowed(WriteMethodNotAllowedResponse)

	return &h
}
