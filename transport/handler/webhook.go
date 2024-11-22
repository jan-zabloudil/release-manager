package handler

import (
	"errors"
	"io"
	"net/http"

	resperrors "release-manager/transport/errors"
	"release-manager/transport/model"
	"release-manager/transport/util"
)

const (
	githubWebhookDeleteEventName = "delete"
	SignatureHeader              = "X-Hub-Signature-256"
	// GithubHookEvent is the header key for the GitHub webhook event type.
	// Docs: https://docs.github.com/en/webhooks/webhook-events-and-payloads
	GithubHookEvent = "X-GitHub-Event"
)

func (h *Handler) handleGithubTagDeletionWebhook(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get(GithubHookEvent)
	if event != githubWebhookDeleteEventName {
		util.WriteResponseError(w, resperrors.NewDefaultBadRequestError().Wrap(errors.New("not a delete event")))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteResponseError(w, resperrors.NewInvalidRequestPayloadError().Wrap(err))
		return
	}

	input := model.ToSvcGithubTagDeletionWebhookInput(
		body,
		r.Header.Get(SignatureHeader),
	)

	if err := h.ReleaseSvc.DeleteReleaseOnGitTagRemoval(
		r.Context(),
		input,
	); err != nil {
		util.WriteResponseError(w, resperrors.NewFromSvcErr(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
