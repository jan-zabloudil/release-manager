SELECT
    e.id,
    e.project_id,
    e.name,
    e.service_url,
    e.created_at,
    e.updated_at,
    r.id AS release_id,
    r.release_title,
    r.release_notes,
    r.created_by AS release_created_by,
    r.created_at AS release_created_at,
    r.updated_at AS release_updated_at
FROM environments e
LEFT JOIN releases r
    ON r.id = e.deployed_release_id
WHERE
    e.project_id = @projectID AND
    e.name = @name
