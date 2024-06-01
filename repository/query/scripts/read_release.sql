SELECT
    r.*,
    p.github_owner_slug,
    p.github_repo_slug
FROM releases r
JOIN projects p
    ON r.project_id = p.id
WHERE
    r.id = @releaseID AND
    r.project_id = @projectID
