package transport

import (
	"net/http"
)

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	err := WriteResponse(w, http.StatusNoContent, nil, nil)
	if err != nil {
		WriteServerErrorResponse(w, r, err)
	}
}
