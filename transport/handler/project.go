package handler

import (
	"net/http"

	"release-manager/pkg/id"
	resperr "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var input model.CreateProjectInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewFromBodyUnmarshalErr(err))
		return
	}

	p, err := h.ProjectSvc.CreateProject(
		r.Context(),
		model.ToSvcCreateProjectInput(input),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToProject(p))
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.GetProject(r.Context(), projectID, util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProject(p))
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.ListProjects(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToProjects(p))
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.UpdateProjectInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.UpdateProject(
		r.Context(),
		model.ToSvcUpdateProjectInput(input),
		projectID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.DeleteProject(r.Context(), projectID, util.ContextAuthUserID(r)); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listGithubRepoTags(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	t, err := h.ProjectSvc.ListGithubRepoTags(
		r.Context(),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToGitTags(t))
}

func (h *Handler) setGithubRepoForProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.SetProjectGithubRepoInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ProjectSvc.SetGithubRepoForProject(
		r.Context(),
		input.RawRepoURL,
		projectID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getGithubRepoForProject(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperr.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	repo, err := h.ProjectSvc.GetGithubRepoForProject(
		r.Context(),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperr.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToGithubRepo(repo))
}
