package model

type GithubRefWebhookInput struct {
	// Ref is the name of the reference (tag, branch, etc.)
	Ref string `json:"ref" required:"true"`
	// RefType is the type of the reference (e.g. "branch" or "tag")
	RefType string `json:"ref_type" required:"true"`
	Repo    struct {
		// Owner and repository slug of the GitHub repository separated by a slash
		// (e.g. "owner/repo")
		Slugs string `json:"full_name"`
	} `json:"repository"`
}
