package handler

import (
	"net/http"

	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) createRelease(w http.ResponseWriter, r *http.Request) {
	var input model.Release
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	rls, err := model.NewSvcRelease(
		utils.ContextApp(r).ID,
		input.SourceCode,
		input.Title,
		input.ChangeLog,
		input.Deployments.Dev,
		input.Deployments.Stg,
		input.Deployments.Prd,
		utils.ContextUser(r).ID,
	)

	if err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	rls, err = h.ReleaseSvc.Create(r.Context(), rls)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, model.ToNetRelease(
		rls.ID,
		rls.SourceCode.Tag(),
		rls.SourceCode.TargetCommitIsh(),
		rls.Deployments.Dev,
		rls.Deployments.Stg,
		rls.Deployments.Prd,
		rls.Title,
		rls.ChangeLog,
		rls.CreatedByUserID,
		rls.CreatedAt,
		rls.UpdatedAt,
	))
}

func (h *Handler) listReleases(w http.ResponseWriter, r *http.Request) {
	releases, err := h.ReleaseSvc.GetAllForApp(r.Context(), utils.ContextApp(r).ID)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetReleases(releases))
}

func (h *Handler) getRelease(w http.ResponseWriter, r *http.Request) {
	rls := utils.ContextRelease(r)
	utils.WriteJSONResponse(w, http.StatusCreated, model.ToNetRelease(
		rls.ID,
		rls.SourceCode.Tag(),
		rls.SourceCode.TargetCommitIsh(),
		rls.Deployments.Dev,
		rls.Deployments.Stg,
		rls.Deployments.Prd,
		rls.Title,
		rls.ChangeLog,
		rls.CreatedByUserID,
		rls.CreatedAt,
		rls.UpdatedAt,
	))
}

func (h *Handler) deleteRelease(w http.ResponseWriter, r *http.Request) {
	err := h.ReleaseSvc.Delete(r.Context(), utils.ContextRelease(r).ID)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateRelease(w http.ResponseWriter, r *http.Request) {

	var input model.ReleasePatch
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	rls, err := model.ToSvcRelease(
		*utils.ContextRelease(r),
		input.SourceCode,
		input.Title,
		input.ChangeLog,
		input.Deployments.Dev,
		input.Deployments.Stg,
		input.Deployments.Prd,
	)
	if err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	rls, err = h.ReleaseSvc.Update(r.Context(), rls)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetRelease(
		rls.ID,
		rls.SourceCode.Tag(),
		rls.SourceCode.TargetCommitIsh(),
		rls.Deployments.Dev,
		rls.Deployments.Stg,
		rls.Deployments.Prd,
		rls.Title,
		rls.ChangeLog,
		rls.CreatedByUserID,
		rls.CreatedAt,
		rls.UpdatedAt,
	))
}
