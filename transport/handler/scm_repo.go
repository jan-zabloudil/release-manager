package handler

import (
	"errors"
	"net/http"

	githuberr "release-manager/github/errors"
	svcerr "release-manager/service/errors"
	"release-manager/transport/model"
	"release-manager/transport/utils"

	httpx "go.strv.io/net/http"
)

func (h *Handler) setSCMRepo(w http.ResponseWriter, r *http.Request) {
	var input model.SCMRepo
	if err := utils.UnmarshalRequest(r, &input); err != nil {
		utils.WriteBadRequestResponse(w, err)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	repo, err := model.ToSvcSCMRepo(utils.ContextApp(r).ID, input.Platform, input.RepoURL)
	if err != nil {
		utils.WriteUnprocessableEntityResponse(w, err)
		return
	}

	repo, err = h.SCMRepoSvc.SetRepo(r.Context(), repo)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, model.ToNetSCMRepo(repo.Platform(), repo.RepoAbsURL()))
}

func (h *Handler) getSCMRepo(w http.ResponseWriter, r *http.Request) {
	repo, err := h.SCMRepoSvc.GetRepo(r.Context(), utils.ContextApp(r).ID)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetSCMRepo(repo.Platform(), repo.RepoAbsURL()))
}

func (h *Handler) deleteSCMRepo(w http.ResponseWriter, r *http.Request) {
	err := h.SCMRepoSvc.DeleteRepo(r.Context(), utils.ContextApp(r).ID)
	if err != nil {
		utils.WriteServerErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getSCMRepoTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.SCMRepoSvc.GetTags(r.Context(), utils.ContextApp(r).ID)
	if err != nil {
		switch {
		case errors.Is(err, githuberr.ErrResourceNotFound):
			utils.WriteErrorResponse(
				w,
				http.StatusNotFound,
				httpx.WithError(err),
				httpx.WithErrorCode("404"),
				httpx.WithErrorMessage(err.Error()),
			)
			return
		case errors.Is(err, githuberr.ErrUnauthenticated):
			utils.WriteErrorResponse(
				w,
				http.StatusUnauthorized,
				httpx.WithError(err),
				httpx.WithErrorCode("401"),
				httpx.WithErrorMessage(err.Error()),
			)
			return
		case errors.Is(err, githuberr.ErrForbidden):
			utils.WriteErrorResponse(
				w,
				http.StatusForbidden,
				httpx.WithError(err),
				httpx.WithErrorCode("403"),
				httpx.WithErrorMessage(err.Error()),
			)
			return
		case errors.Is(err, svcerr.ErrSCMRepoNotSet):
			utils.WriteConflictResponse(w, err)
			return
		default:
			utils.WriteServerErrorResponse(w, err)
			return
		}
	}

	utils.WriteJSONResponse(w, http.StatusOK, model.ToNetGitTags(tags))
}
