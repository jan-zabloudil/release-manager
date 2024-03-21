package transport

import (
	"net/http"

	neterrs "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
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

	h.Mux.Use(httpx.RequestIDMiddleware(RequestID))
	h.Mux.Use(httpx.LoggingMiddleware(utils.NewServerLogger("logging")))
	h.Mux.Use(httpx.RecoverMiddleware(utils.NewServerLogger("recover")))
	h.Mux.Use(h.auth)

	h.Mux.Get("/admin/users", h.requireAdminUser(h.getUsers))
	h.Mux.Route("/admin/users/{id}", func(r chi.Router) {
		r.Get("/", h.requireAdminUser(h.handleUser))
		r.Delete("/", h.requireAdminUser(h.handleUser))
	})

	h.Mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.Mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		WriteNotFoundResponse(w, neterrs.ErrHttpNotFound)
	})
	h.Mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		WriteMethodNotAllowedResponse(w, neterrs.ErrHttpMethodNotAllowed)
	})

	return &h
}
