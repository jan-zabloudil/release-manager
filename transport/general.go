package transport

import (
	"net/http"

	neterr "github.com/jan-zabloudil/release-manager/transport/errors"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	err := WriteResponse(w, http.StatusNoContent, nil, nil)
	if err != nil {
		WriteServerErrorResponse(w, r, err)
	}
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	WriteNotFoundResponse(w, r, neterr.ErrHttpNotFound)
}

func (h *Handler) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	WriteMethodNotAllowedResponse(w, r, neterr.ErrHttpMethodNotAllowed)
}
