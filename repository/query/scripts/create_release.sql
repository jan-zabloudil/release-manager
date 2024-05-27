INSERT INTO releases (
    id,
    project_id,
    release_title,
    release_notes,
    created_by,
    created_at,
    updated_at,
    git_tag_name,
    github_release_created_at,
    github_release_updated_at
)
VALUES (
    @id,
    @projectID,
    @releaseTitle,
    @releaseNotes,
    @createdBy,
    @createdAt,
    @updatedAt,
    @gitTagName,
    @githubReleaseCreatedAt,
    @githubReleaseUpdatedAt
)
