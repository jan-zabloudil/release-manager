package handler

import (
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProjectInput
	if err := util.UnmarshalRequest(r, &req); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.CreateProject(
		r.Context(),
		model.ToSvcCreateProjectInput(req),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToProject(p))
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.GetProject(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProject(p))
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.ListProjects(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjects(p))
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateProjectInput

	if err := util.UnmarshalRequest(r, &req); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.UpdateProject(
		r.Context(),
		model.ToSvcUpdateProjectInput(req),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProject(p))
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := h.ProjectSvc.DeleteProject(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r)); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listGithubRepositoryTags(w http.ResponseWriter, r *http.Request) {
	t, err := h.ProjectSvc.ListGithubRepositoryTags(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToGitTags(t))
}
