package transport

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	svcerr "release-manager/service/errors"
	svcmodel "release-manager/service/model"
	"release-manager/transport/model"
	"release-manager/transport/utils"

	"github.com/google/uuid"
)

func (h *Handler) listProjectMembers(w http.ResponseWriter, r *http.Request) {
	m, err := h.ProjectMemberSvc.ListAll(r.Context(), ContextProject(r).ID)
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetProjectMembers(m))
}

func (h *Handler) handleProjectMember(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUUIDParamFrom(r, "userId")
	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	m, err := h.ProjectMemberSvc.Get(r.Context(), ContextProject(r).ID, userID)
	if err != nil {
		switch {
		case errors.Is(err, reperr.ErrResourceNotFound):
			WriteNotFoundResponse(w, err)
			return
		default:
			WriteServerErrorResponse(w, err)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.getProjectMember(w, model.ToNetProjectMember(m))
		return
	case http.MethodPatch:
		h.updateProjectMember(w, r, m)
	case http.MethodDelete:
		h.deleteProjectMember(w, r, userID)
	}
}

func (h *Handler) getProjectMember(w http.ResponseWriter, member model.ProjectMember) {
	WriteJSONResponse(w, http.StatusOK, member)
}

func (h *Handler) deleteProjectMember(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	if err := h.ProjectMemberSvc.Delete(r.Context(), ContextProject(r).ID, userID); err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateProjectMember(w http.ResponseWriter, r *http.Request, m svcmodel.ProjectMember) {
	var input model.UpdateProjectRole

	if err := UnmarshalRequest(r, &input); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		WriteUnprocessableEntityResponse(w, err)
		return
	}

	role, err := svcmodel.NewProjectRole(input.Role)
	if err != nil {
		WriteUnprocessableEntityResponse(w, err)
		return
	}

	m, err = h.ProjectMemberSvc.UpdateRole(r.Context(), m, ContextProjectMember(r), role)
	if err != nil {
		switch {
		case errors.Is(err, svcerr.ErrProjectMemberRoleCannotBeGranted), errors.Is(err, svcerr.ErrProjectMemberUpdateNotAllowed):
			WriteForbiddenErrorResponse(w, err)
			return
		default:
			WriteServerErrorResponse(w, err)
			return
		}
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetProjectMember(m))
}
