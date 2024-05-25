package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyGithubWebhookPayload verifies the payload of a GitHub webhook
// using the secret and the signature provided in X-Hub-Signature-256 header
// Docs: https://docs.github.com/en/webhooks/using-webhooks/validating-webhook-deliveries
func VerifyGithubWebhookPayload(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	// always returns nil error
	_, _ = mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
