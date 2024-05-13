package query

import (
	_ "embed"
)

var (
	//go:embed scripts/create_release.sql
	CreateRelease string
	//go:embed scripts/read_release_for_project.sql
	ReadReleaseForProject string
)
