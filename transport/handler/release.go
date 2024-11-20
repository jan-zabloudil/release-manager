package handler

import (
	"net/http"

	"release-manager/pkg/id"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createRelease(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.CreateReleaseInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	rls, err := h.ReleaseSvc.CreateRelease(
		r.Context(),
		model.ToSvcCreateReleaseInput(input),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToRelease(rls))
}

func (h *Handler) getRelease(w http.ResponseWriter, r *http.Request) {
	rlsID, err := util.GetPathParam[id.Release](r, "release_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	rls, err := h.ReleaseSvc.GetRelease(
		r.Context(),
		rlsID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToRelease(rls))
}

func (h *Handler) deleteRelease(w http.ResponseWriter, r *http.Request) {
	rlsID, err := util.GetPathParam[id.Release](r, "release_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.DeleteReleaseInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ReleaseSvc.DeleteRelease(
		r.Context(),
		model.ToSvcDeleteReleaseInput(input),
		rlsID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listReleases(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	rls, err := h.ReleaseSvc.ListReleasesForProject(
		r.Context(),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToReleases(rls))
}

func (h *Handler) updateRelease(w http.ResponseWriter, r *http.Request) {
	rlsID, err := util.GetPathParam[id.Release](r, "release_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.UpdateReleaseInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ReleaseSvc.UpdateRelease(
		r.Context(),
		model.ToSvcUpdateReleaseInput(input),
		rlsID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) sendReleaseNotification(w http.ResponseWriter, r *http.Request) {
	rlsID, err := util.GetPathParam[id.Release](r, "release_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ReleaseSvc.SendReleaseNotification(
		r.Context(),
		rlsID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) upsertGithubRelease(w http.ResponseWriter, r *http.Request) {
	rlsID, err := util.GetPathParam[id.Release](r, "release_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ReleaseSvc.UpsertGithubRelease(
		r.Context(),
		rlsID,
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) generateGithubReleaseNotes(w http.ResponseWriter, r *http.Request) {
	projectID, err := util.GetPathParam[id.Project](r, "project_id")
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	var input model.GithubReleaseNotesInput
	if err := util.UnmarshalBody(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	notes, err := h.ReleaseSvc.GenerateGithubReleaseNotes(
		r.Context(),
		model.ToSvcGithubReleaseNotesInput(input),
		projectID,
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToGithubReleaseNotes(notes))
}
