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

	rls, err := h.ReleaseSvc.Create(
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
	rls, err := h.ReleaseSvc.Get(
		r.Context(),
		util.ContextProjectID(r),
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
	if err := h.ReleaseSvc.Delete(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listReleases(w http.ResponseWriter, r *http.Request) {
	rls, err := h.ReleaseSvc.ListForProject(
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

	rls, err := h.ReleaseSvc.Update(
		r.Context(),
		model.ToSvcUpdateReleaseInput(input),
		util.ContextProjectID(r),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, model.ToRelease(rls))
}

func (h *Handler) sendReleaseNotification(w http.ResponseWriter, r *http.Request) {
	if err := h.ReleaseSvc.SendReleaseNotification(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextReleaseID(r),
		util.ContextAuthUserID(r),
	); err != nil {
		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
