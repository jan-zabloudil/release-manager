package handler

import (
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createRelease(w http.ResponseWriter, r *http.Request) {
	var input model.CreateReleaseInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
		util.WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	rls, err := h.ReleaseSvc.Create(
		r.Context(),
		model.ToSvcCreateReleaseInput(input),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, util.ToResponseError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToRelease(rls))
}