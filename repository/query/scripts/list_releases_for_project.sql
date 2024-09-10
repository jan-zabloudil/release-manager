SELECT
    r.*,
    p.github_owner_slug,
    p.github_repo_slug,
    COALESCE(JSON_AGG(
        JSON_BUILD_OBJECT(
                'id', ra.attachment_id,
                'name', ra.name,
                'file_path', ra.file_path,
                'created_at', ra.created_at
        )
    ) FILTER (WHERE ra.release_id IS NOT NULL), '[]') AS attachments
FROM releases r
JOIN projects p
    ON r.project_id = p.id
LEFT JOIN release_attachments ra
    ON ra.release_id = r.id
WHERE r.project_id = @projectID
GROUP BY r.id, p.github_owner_slug, p.github_repo_slug
ORDER BY r.created_at DESC
