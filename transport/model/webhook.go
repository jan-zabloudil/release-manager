package model

type GithubReleaseWebhookInput struct {
	Action  string `json:"action" required:"true"`
	Release struct {
		TagName string `json:"tag_name" required:"true"`
		Name    string `json:"name" required:"true"`
		Body    string `json:"body" required:"true"`
	} `json:"release"`
	Repo struct {
		// Owner and repository slug of the GitHub repository separated by a slash
		// (e.g. "owner/repo")
		Slugs string `json:"full_name" required:"true"`
	} `json:"repository"`
}
