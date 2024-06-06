package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"release-manager/pkg/validator"
	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

const (
	githubWebhookDeleteEventName = "delete"
	gitRefTypeTag                = "tag"
)

func (h *Handler) handleGithubReleaseWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err))
		return
	}

	secret, err := h.SettingsSvc.GetGithubWebhookSecret(r.Context())
	if err != nil {
		if resperrors.IsGithubIntegrationNotEnabledError(err) {
			// If the GitHub integration is not enabled, the webhook should not be processed.
			// And we should return 204, because from the webhook perspective it is not an error case.
			w.WriteHeader(http.StatusNoContent)
			return
		}

		util.WriteResponseError(w, resperrors.ToError(err))
		return
	}

	if !isValidPayload(secret, r.Header.Get(util.SignatureHeader), body) {
		util.WriteResponseError(w, resperrors.NewUnauthorizedError())
		return
	}

	var input model.GithubRefWebhookInput
	if err := json.Unmarshal(body, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err))
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err))
		return
	}

	if r.Header.Get(util.GithubHookEvent) == githubWebhookDeleteEventName && input.RefType == gitRefTypeTag {
		ownerSlug, repoSlug, err := model.SplitGithubRepoSlugs(input.Repo.Slugs)
		if err != nil {
			util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err))
			return
		}

		err = h.ReleaseSvc.DeleteReleaseByGitTag(r.Context(), ownerSlug, repoSlug, input.Ref)
		if err != nil && !resperrors.IsNotFoundError(err) {
			// If the release is not found, from the webhook perspective that is not an error.
			// Webhook just notifies us that the tag is deleted but there is no release for that tag.
			// Therefore, we want to return error only if the error is not NotFoundError.
			util.WriteResponseError(w, resperrors.ToError(err))
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func isValidPayload(secret, signature string, payload []byte) bool {
	// if secret is not set, the payload is not verified
	if secret == "" {
		return true
	}

	// if secret is set, the service requires a webhook that provides a signature to verify the payload
	if signature == "" {
		return false
	}

	if !util.VerifyGithubWebhookPayload(secret, signature, payload) {
		return false
	}

	return true
}
