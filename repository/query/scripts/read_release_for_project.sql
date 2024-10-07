-- The WITH clause is used to pre-aggregate the attachments, allowing us to avoid
-- a GROUP BY in the main query. If we joined release_attachments directly with
-- releases, a GROUP BY would be required, making it impossible to use FOR UPDATE
-- due to its incompatibility with GROUP BY.
-- FOR UPDATE is appended in the repository function when reading the release
-- while updating it.
WITH attachments AS (
    SELECT
        ra.release_id,
        JSON_AGG(
                JSON_BUILD_OBJECT(
                        'attachment_id', ra.attachment_id,
                        'name', ra.name,
                        'file_path', ra.file_path,
                        'created_at', ra.created_at
                )
        ) AS attachments
    FROM release_attachments ra
    GROUP BY ra.release_id
)
SELECT
    r.*,
    p.github_owner_slug,
    p.github_repo_slug,
    COALESCE(a.attachments, '[]'::json) AS attachments
FROM releases r
JOIN projects p
    ON r.project_id = p.id
LEFT JOIN attachments a
    ON a.release_id = r.id
WHERE
    r.id = @releaseID AND
    r.project_id = @projectID
