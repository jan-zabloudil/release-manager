package handler

import (
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createDeployment(w http.ResponseWriter, r *http.Request) {
	var input model.CreateDeploymentInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	dpl, err := h.ReleaseSvc.CreateDeployment(
		r.Context(),
		model.ToSvcCreateDeploymentInput(input),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToDeployment(dpl))
}

func (h *Handler) listDeploymentsForProject(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.ListDeploymentsFilterParams](r)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	dpls, err := h.ReleaseSvc.ListDeploymentsForProject(
		r.Context(),
		model.ToSvcListDeploymentsFilterParams(params),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToDeployments(dpls))
}
