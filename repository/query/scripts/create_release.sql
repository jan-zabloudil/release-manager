INSERT INTO releases (
    id,
    project_id,
    release_title,
    release_notes,
    created_by,
    created_at,
    updated_at,
    github_release_id,
    github_owner_slug,
    github_repo_slug,
    github_release_data
)
VALUES (
    @id,
    @projectID,
    @releaseTitle,
    @releaseNotes,
    @createdBy,
    @createdAt,
    @updatedAt,
    @githubReleaseID,
    @githubOwnerSlug,
    @githubRepoSlug,
    @githubReleaseData
)
