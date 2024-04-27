package transport

import (
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProjectInput
	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.Create(
		r.Context(),
		model.ToSvcCreateProjectInput(req),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusCreated, model.ToProject(p))
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.Get(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r))
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToProject(p))
}

func (h *Handler) listProjects(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.ListAll(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToProjects(p))
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateProjectInput

	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.Update(
		r.Context(),
		model.ToSvcUpdateProjectInput(req),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToProject(p))
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := h.ProjectSvc.Delete(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r)); err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
