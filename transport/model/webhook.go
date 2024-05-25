package model

type GithubReleaseWebhookInput struct {
	Action  string `json:"action" required:"true"`
	Release struct {
		Name string `json:"name" required:"true"`
		Body string `json:"body"`
	} `json:"release"`
	Repo struct {
		// Owner and repository slug of the GitHub repository separated by a slash
		// (e.g. "owner/repo")
		Slugs string `json:"full_name"`
	} `json:"repository"`
}
