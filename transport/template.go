package transport

import (
	"errors"
	"net/http"

	reperr "release-manager/repository/errors"
	svcmodel "release-manager/service/model"
	"release-manager/transport/model"

	"github.com/go-playground/validator/v10"
)

func (h *Handler) createTemplate(w http.ResponseWriter, r *http.Request) {
	var input model.Template
	if err := UnmarshalRequest(r, &input); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			WriteUnprocessableEntityResponse(w, err)
			return
		}

		WriteBadRequestResponse(w, err)
		return
	}

	t, err := model.NewSvcTemplate(input.Type, model.NewSvcReleaseMsg(input.ReleaseMsg.Title, input.ReleaseMsg.Text, input.ReleaseMsg.Includes))
	if err != nil {
		WriteUnprocessableEntityResponse(w, err)
		return
	}

	t, err = h.TemplateSvc.Create(r.Context(), t)
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(
		w,
		http.StatusCreated,
		model.ToNetTemplate(
			t.ID,
			t.Type.TemplateType(),
			model.ToNetReleaseMsg(t.ReleaseMsg.Title, t.ReleaseMsg.Text, t.ReleaseMsg.Includes),
		),
	)
}

func (h *Handler) listTemplates(w http.ResponseWriter, r *http.Request) {
	t, err := h.TemplateSvc.GetAll(r.Context())
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusOK, model.ToNetTemplates(t))
}

func (h *Handler) handleTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromURL(r)
	if err != nil {
		WriteNotFoundResponse(w, err)
		return
	}

	t, err := h.TemplateSvc.Get(r.Context(), id)
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
		h.getTemplate(w, t)
	case http.MethodPatch:
		h.updateTemplate(w, r, t)
	case http.MethodDelete:
		h.deleteTemplate(w, r, t)
	}
}

func (h *Handler) updateTemplate(w http.ResponseWriter, r *http.Request, t svcmodel.Template) {
	var input model.ReleaseMessagePatch
	if err := UnmarshalRequest(r, &input); err != nil {
		WriteBadRequestResponse(w, err)
		return
	}

	t, err := h.TemplateSvc.Update(
		r.Context(),
		model.ToSvcTemplate(t, model.ToSvcReleaseMsg(t.ReleaseMsg, input.Title, input.Text, input.Includes)),
	)
	if err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	WriteJSONResponse(w, http.StatusCreated, model.ToNetTemplate(
		t.ID,
		t.Type.TemplateType(),
		model.ToNetReleaseMsg(t.ReleaseMsg.Title, t.ReleaseMsg.Text, t.ReleaseMsg.Includes)),
	)
}

func (h *Handler) getTemplate(w http.ResponseWriter, t svcmodel.Template) {
	WriteJSONResponse(w, http.StatusOK, model.ToNetTemplate(
		t.ID,
		t.Type.TemplateType(),
		model.ToNetReleaseMsg(t.ReleaseMsg.Title, t.ReleaseMsg.Text, t.ReleaseMsg.Includes)),
	)
}

func (h *Handler) deleteTemplate(w http.ResponseWriter, r *http.Request, t svcmodel.Template) {
	if err := h.TemplateSvc.Delete(r.Context(), t.ID); err != nil {
		WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
