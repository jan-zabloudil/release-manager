package handler

import (
	"encoding/json"
	"io"
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

const (
	githubReleaseWebhookEditedAction  = "edited"
	githubReleaseWebhookDeletedAction = "deleted"
)

func (h *Handler) handleGithubReleaseWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err))
		return
	}

	// TODO get the secret from the settings
	secret := ""

	if !isValidPayload(secret, r.Header.Get(util.SignatureHeader), body) {
		util.WriteResponseError(w, resperrors.NewUnauthorizedError())
		return
	}

	var input model.GithubReleaseWebhookInput
	if err := json.Unmarshal(body, &input); err != nil {
		util.WriteResponseError(w, resperrors.NewBadRequestError().Wrap(err))
		return
	}

	switch input.Action {
	case githubReleaseWebhookEditedAction:
		// TODO call the service to update the release
	case githubReleaseWebhookDeletedAction:
		// TODO call the service to delete the release
	}

	w.WriteHeader(http.StatusOK)
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
