package model

import (
	svcmodel "release-manager/service/model"
)

func ToSvcGithubTagDeletionWebhookInput(payload []byte, signature string) svcmodel.GithubTagDeletionWebhookInput {
	return svcmodel.GithubTagDeletionWebhookInput{
		RawPayload: payload,
		Signature:  signature,
	}
}
