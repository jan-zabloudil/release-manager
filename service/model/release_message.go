package model

type ReleaseMessage struct {
	Title    string
	Text     string
	Includes Includes
}

// Includes control which properties are included in release message
type Includes struct {
	ProjectName   bool
	AppName       bool
	ReleaseName   bool
	Changelog     bool
	Deployments   bool
	GithubRelease bool
	GithubTag     bool
}
