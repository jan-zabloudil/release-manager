package handler

import (
	"net/http"

	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var input model.Project
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	p, err := h.ProjectSvc.Create(r.Context(), model.ToSvcProject(input), utils.ContextUser(r).ID)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, model.ToNetProject(p))
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, utils.ContextProject(r))
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var input model.ProjectPatch

	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}
	netp := model.PatchToNetProject(input, *utils.ContextProject(r))

	if err := utils.Validate.Struct(netp); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	p, err := h.ProjectSvc.Update(r.Context(), model.ToSvcProject(netp))
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetProject(p))
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.ListAll(r.Context(), model.ToSvcUser(*utils.ContextUser(r)))
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetProjects(p))
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	p := utils.ContextProject(r)
	if err := h.ProjectSvc.Delete(r.Context(), p.ID); err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
