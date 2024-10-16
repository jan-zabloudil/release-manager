package handler

import (
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createRelease(w http.ResponseWriter, r *http.Request) {
	var input model.CreateReleaseInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	rls, err := h.ReleaseSvc.CreateRelease(
		r.Context(),
		model.ToSvcCreateReleaseInput(input),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.ToRelease(rls))
}

func (h *Handler) getRelease(w http.ResponseWriter, r *http.Request) {
	rls, err := h.ReleaseSvc.GetRelease(
		r.Context(),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToRelease(rls))
}

func (h *Handler) deleteRelease(w http.ResponseWriter, r *http.Request) {
	var input model.DeleteReleaseInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ReleaseSvc.DeleteRelease(
		r.Context(),
		model.ToSvcDeleteReleaseInput(input),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listReleases(w http.ResponseWriter, r *http.Request) {
	rls, err := h.ReleaseSvc.ListReleasesForProject(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToReleases(rls))
}

func (h *Handler) updateRelease(w http.ResponseWriter, r *http.Request) {
	var input model.UpdateReleaseInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	if err := h.ReleaseSvc.UpdateRelease(
		r.Context(),
		model.ToSvcUpdateReleaseInput(input),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) sendReleaseNotification(w http.ResponseWriter, r *http.Request) {
	if err := h.ReleaseSvc.SendReleaseNotification(
		r.Context(),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) upsertGithubRelease(w http.ResponseWriter, r *http.Request) {
	if err := h.ReleaseSvc.UpsertGithubRelease(
		r.Context(),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) generateGithubReleaseNotes(w http.ResponseWriter, r *http.Request) {
	var input model.GithubGeneratedReleaseNotesInput
	if err := util.UnmarshalRequest(r, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	notes, err := h.ReleaseSvc.GenerateGithubReleaseNotes(
		r.Context(),
		model.ToSvcGithubGeneratedReleaseNotesInput(input),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToGithubGeneratedReleaseNotes(notes))
}
