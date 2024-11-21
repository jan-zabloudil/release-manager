package handler

import (
	"net/http"

	"release-manager/pkg/id"
	resperr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createDeployment(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage("Invalid project ID"))
		return
	}

	var input model.CreateDeploymentInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	dpl, err := h.ReleaseSvc.CreateDeployment(
		r.Context(),
		model.ToSvcCreateDeploymentInput(input),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToDeployment(dpl))
}

func (h *Handler) listDeploymentsForProject(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.ListDeploymentsParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	dpls, err := h.ReleaseSvc.ListDeploymentsForProject(
		r.Context(),
		model.ToSvcListDeploymentsFilterParams(params),
		params.ProjectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToDeployments(dpls))
}
