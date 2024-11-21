package handler

import (
	"net/http"

	"release-manager/pkg/id"
	resperr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createEnvironment(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.CreateEnvironmentInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	env, err := h.ProjectSvc.CreateEnvironment(
		r.Context(),
		model.ToSvcCreateEnvironmentInput(input, projectID),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToEnvironment(env))
}

func (h *Handler) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.EnvironmentURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.UpdateEnvironmentInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	if err := h.ProjectSvc.UpdateEnvironment(
		r.Context(),
		model.ToSvcUpdateEnvironmentInput(input),
		params.ProjectID,
		params.EnvironmentID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getEnvironment(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.EnvironmentURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	env, err := h.ProjectSvc.GetEnvironment(
		r.Context(),
		params.ProjectID,
		params.EnvironmentID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToEnvironment(env))
}

func (h *Handler) listEnvironments(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	envs, err := h.ProjectSvc.ListEnvironments(
		r.Context(),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToEnvironments(envs))
}

func (h *Handler) deleteEnvironment(w http.ResponseWriter, r *http.Request) {
	params, err := util.UnmarshalURLParams[model.EnvironmentURLParams](r)
	if err != nil {
		util.WriteResponseError(w, resperr.NewInvalidURLParamsError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.DeleteEnvironment(
		r.Context(),
		params.ProjectID,
		params.EnvironmentID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
