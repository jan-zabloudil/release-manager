package transport

import (
	"net/http"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	if err := WriteResponse(w, http.StatusNoContent, nil, nil); err != nil {
		WriteServerErrorResponse(w, r, err)
	}
}
