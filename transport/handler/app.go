package handler

import (
	"net/http"

	"release-manager/transport/model"
	"release-manager/transport/utils"
)

func (h *Handler) createApp(w http.ResponseWriter, r *http.Request) {
	var input model.App
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	app, err := model.NewSvcApp(
		utils.ContextProjectID(r),
		input.Name,
		input.Description,
		input.Environments,
	)
	if err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	app, err = h.AppSvc.Create(r.Context(), app)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, model.ToNetApp(
		app.ID,
		app.Name,
		app.Description,
		app.Environments,
		app.CreatedAt,
		app.UpdatedAt,
	))
}

func (h *Handler) listApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.AppSvc.GetAllForProject(r.Context(), utils.ContextProjectID(r))
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetApps(apps))
}

func (h *Handler) getApp(w http.ResponseWriter, r *http.Request) {
	app := utils.ContextApp(r)
	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetApp(
		app.ID,
		app.Name,
		app.Description,
		app.Environments,
		app.CreatedAt,
		app.UpdatedAt,
	))
}

func (h *Handler) updateApp(w http.ResponseWriter, r *http.Request) {
	var input model.AppPatch
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	app, err := model.ToSvcApp(
		*utils.ContextApp(r),
		input.Name,
		input.Description,
		input.Environments,
	)
	if err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	app, err = h.AppSvc.Update(r.Context(), app)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetApp(
		app.ID,
		app.Name,
		app.Description,
		app.Environments,
		app.CreatedAt,
		app.UpdatedAt,
	))
}

func (h *Handler) deleteApp(w http.ResponseWriter, r *http.Request) {
	app := utils.ContextApp(r)
	if err := h.AppSvc.Delete(r.Context(), app.ID); err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
