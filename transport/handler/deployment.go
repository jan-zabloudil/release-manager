package handler

import (
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createDeployment(w http.ResponseWriter, r *http.Request) {
	var input model.CreateDeploymentInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
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
	filterParams, err := model.ToSvcDeploymentFilterParams(
		util.GetQueryParam(r, "release_id"),
		util.GetQueryParam(r, "environment_id"),
		util.GetQueryParam(r, "latest_only"),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	dpls, err := h.ReleaseSvc.ListDeploymentsForProject(
		r.Context(),
		filterParams,
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToDeployments(dpls))
}
