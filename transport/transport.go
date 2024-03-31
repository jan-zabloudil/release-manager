package transport

import (
	"log/slog"
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"

	"github.com/go-chi/chi/v5"
	httpx "go.strv.io/net/http"
)

type Handler struct {
	Mux     *chi.Mux
	AuthSvc model.AuthService
	UserSvc model.UserService
}

func NewHandler(as model.AuthService, us model.UserService) *Handler {
	h := Handler{
		Mux:     chi.NewRouter(),
		AuthSvc: as,
		UserSvc: us,
	}

	h.Mux.Use(httpx.RequestIDMiddleware(RequestID))
	h.Mux.Use(httpx.LoggingMiddleware(slog.Default().WithGroup("logger")))
	h.Mux.Use(httpx.RecoverMiddleware(slog.Default().WithGroup("recover")))
	h.Mux.Use(h.auth)

	h.Mux.Get("/admin/users", h.requireAuthUser(h.getUsers))
	h.Mux.Route("/admin/users/{id}", func(r chi.Router) {
		r.Use(h.handleResourceID("id", util.ContextSetUserID))
		r.Get("/", h.requireAuthUser(h.getUser))
		r.Delete("/", h.requireAuthUser(h.deleteUser))
	})

	h.Mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.Mux.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		WriteResponseError(w, responseerrors.NewNotFoundError())
	})
	h.Mux.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		WriteResponseError(w, responseerrors.NewMethodNotAllowedError())
	})

	return &h
}
