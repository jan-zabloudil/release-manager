package transport

import (
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var req model.CreateProjectRequest
	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.Create(
		r.Context(),
		model.ToSvcProjectCreation(
			req.Name,
			req.SlackChannelID,
			req.ReleaseNotificationConfig,
		),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(
		w,
		http.StatusCreated,
		model.ToProjectResponse(
			p.ID,
			p.Name,
			p.SlackChannelID,
			p.ReleaseNotificationConfig,
			p.CreatedAt,
			p.UpdatedAt,
		),
	)
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.Get(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r))
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(
		w,
		http.StatusOK,
		model.ToProjectResponse(
			p.ID,
			p.Name,
			p.SlackChannelID,
			p.ReleaseNotificationConfig,
			p.CreatedAt,
			p.UpdatedAt,
		),
	)
}

func (h *Handler) getProjects(w http.ResponseWriter, r *http.Request) {
	p, err := h.ProjectSvc.GetAll(r.Context(), util.ContextAuthUserID(r))
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToProjects(p))
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateProjectRequest

	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	p, err := h.ProjectSvc.Update(
		r.Context(),
		model.ToSvcProjectUpdate(
			req.Name,
			req.SlackChannelID,
			req.ReleaseNotificationConfig,
		),
		util.ContextProjectID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(
		w,
		http.StatusOK,
		model.ToProjectResponse(
			p.ID,
			p.Name,
			p.SlackChannelID,
			p.ReleaseNotificationConfig,
			p.CreatedAt,
			p.UpdatedAt,
		),
	)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := h.ProjectSvc.Delete(r.Context(), util.ContextProjectID(r), util.ContextAuthUserID(r)); err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
