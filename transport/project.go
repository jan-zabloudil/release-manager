package transport

import (
	"net/http"

	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var input model.Project
	if err := UnmarshalRequest(r, &input); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		WriteUnprocessableEntityResponse(w, err)
		return
	}

	p, err := h.ProjectSvc.Create(r.Context(), model.ToSvcProject(input))
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusCreated, model.ToNetProject(p))
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	WriteJSONResponse(w, http.StatusOK, ContextProject(r))
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var input model.ProjectPatch

	if err := UnmarshalRequest(r, &input); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}
	netp := model.PatchToNetProject(input, *ContextProject(r))

	if err := utils.Validate.Struct(netp); err != nil {
		WriteUnprocessableEntityResponse(w, err)
		return
	}

	p, err := h.ProjectSvc.Update(r.Context(), model.ToSvcProject(netp))
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetProject(p))
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.ListAll(r.Context())
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetProjects(p))
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	p := ContextProject(r)
	if err := h.ProjectSvc.Delete(r.Context(), p.ID); err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
