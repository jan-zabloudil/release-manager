package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jan-zabloudil/release-manager/transport/model"
)

type Handler struct {
	Mux     *chi.Mux
	UserSvc model.UserService
}

func NewHandler(us model.UserService) *Handler {
	h := Handler{
		Mux:     chi.NewRouter(),
		UserSvc: us,
	}
	h.Mux.Use(middleware.Logger) // TODO use slog
	h.Mux.Use(h.recoverPanic)
	h.Mux.Use(h.auth)

	h.Mux.Get("/ping", h.ping)

	h.Mux.NotFound(h.notFound)
	h.Mux.MethodNotAllowed(h.methodNotAllowed)

	return &h
}
