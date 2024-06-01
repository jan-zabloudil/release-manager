SELECT
    r.*,
    p.github_owner_slug,
    p.github_repo_slug
FROM releases r
JOIN projects p
    ON r.project_id = p.id
WHERE r.project_id = @projectID
ORDER BY created_at DESC
