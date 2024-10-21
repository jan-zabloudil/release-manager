package handler

import (
	"net/http"

	"release-manager/pkg/id"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createEnvironment(w http.ResponseWriter, r *http.Request) {
	var req model.CreateEnvironmentInput
	if err := util.UnmarshalRequest(r, &req); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	env, err := h.ProjectSvc.CreateEnvironment(
		r.Context(),
		model.ToSvcCreateEnvironmentInput(req, util.ContextProjectID(r)),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToEnvironment(env))
}

func (h *Handler) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	envID, err := util.GetPathParam[id.Environment](r, "environment_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewInvalidResourceIDError().Wrap(err).WithMessage("invalid environment ID"))
		return
	}

	var req model.UpdateEnvironmentInput
	if err := util.UnmarshalRequest(r, &req); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.UpdateEnvironment(
		r.Context(),
		model.ToSvcUpdateEnvironmentInput(req),
		util.ContextProjectID(r),
		envID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getEnvironment(w http.ResponseWriter, r *http.Request) {
	envID, err := util.GetPathParam[id.Environment](r, "environment_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewInvalidResourceIDError().Wrap(err).WithMessage("invalid environment ID"))
		return
	}

	env, err := h.ProjectSvc.GetEnvironment(
		r.Context(),
		util.ContextProjectID(r),
		envID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToEnvironment(env))
}

func (h *Handler) listEnvironments(w http.ResponseWriter, r *http.Request) {
	envs, err := h.ProjectSvc.ListEnvironments(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToEnvironments(envs))
}

func (h *Handler) deleteEnvironment(w http.ResponseWriter, r *http.Request) {
	envID, err := util.GetPathParam[id.Environment](r, "environment_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewInvalidResourceIDError().Wrap(err).WithMessage("invalid environment ID"))
		return
	}

	if err := h.ProjectSvc.DeleteEnvironment(
		r.Context(),
		util.ContextProjectID(r),
		envID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
