package transport

import (
	"net/http"

	"release-manager/pkg/responseerrors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

func (h *Handler) createEnvironment(w http.ResponseWriter, r *http.Request) {
	var req model.CreateEnvironmentRequest
	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	env, err := h.ProjectSvc.CreateEnvironment(
		r.Context(),
		model.ToSvcEnvironmentCreation(util.ContextProjectID(r), req.Name, req.ServiceURL),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(
		w,
		http.StatusCreated,
		model.ToEnvironmentResponse(env.ID, env.Name, env.ServiceURL, env.CreatedAt, env.UpdatedAt),
	)
}

func (h *Handler) updateEnvironment(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateEnvironmentRequest
	if err := UnmarshalRequest(r, &req); err != nil {
		WriteResponseError(w, responseerrors.NewBadRequestError().Wrap(err).WithMessage(err.Error()))
		return
	}

	env, err := h.ProjectSvc.UpdateEnvironment(
		r.Context(),
		model.ToSvcEnvironmentUpdate(
			req.Name,
			req.ServiceURL,
		),
		util.ContextProjectID(r),
		util.ContextEnvironmentID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(
		w,
		http.StatusOK,
		model.ToEnvironmentResponse(env.ID, env.Name, env.ServiceURL, env.CreatedAt, env.UpdatedAt),
	)
}

func (h *Handler) getEnvironment(w http.ResponseWriter, r *http.Request) {
	env, err := h.ProjectSvc.GetEnvironment(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextEnvironmentID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	WriteJSONResponse(
		w,
		http.StatusOK,
		model.ToEnvironmentResponse(env.ID, env.Name, env.ServiceURL, env.CreatedAt, env.UpdatedAt),
	)
}

func (h *Handler) getEnvironments(w http.ResponseWriter, r *http.Request) {
	envs, err := h.ProjectSvc.GetEnvironments(
		r.Context(),
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
		model.ToEnvironmentsResponse(envs),
	)
}

func (h *Handler) deleteEnvironment(w http.ResponseWriter, r *http.Request) {
	err := h.ProjectSvc.DeleteEnvironment(
		r.Context(),
		util.ContextProjectID(r),
		util.ContextEnvironmentID(r),
		util.ContextAuthUserID(r),
	)
	if err != nil {
		WriteResponseError(w, util.ToResponseError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
