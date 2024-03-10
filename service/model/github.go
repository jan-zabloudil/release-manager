package model

import "context"

type GitHub interface {
	CreateRelease(ctx context.Context, repo string, r GithubRelease) (GithubRelease, error)
	ListTags(ctx context.Context, repo string) ([]GitTag, error)
}

type GithubRelease struct {
	Name            string
	TagName         string
	TargetCommitish string
	Changelog       string
}
