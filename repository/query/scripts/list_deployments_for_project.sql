SELECT
    d.id,
    d.deployed_by,
    d.deployed_at,
    r.id AS release_id,
    r.release_title,
    r.release_notes,
    r.created_by AS release_created_by,
    r.created_at AS release_created_at,
    r.updated_at AS release_updated_at,
    e.id AS env_id,
    e.project_id AS env_project_id,
    e.name AS env_name,
    e.service_url AS env_service_url,
    e.created_at AS env_created_at,
    e.updated_at AS env_updated_at
FROM deployments d
JOIN releases r
    ON d.release_id = r.id
JOIN environments e
    ON d.environment_id = e.id
WHERE
    r.project_id = @projectID AND
    e.project_id = @projectID AND
    (@releaseID::uuid IS NULL OR r.id = @releaseID) AND
    (@envID::uuid IS NULL OR e.id = @envID)
ORDER BY d.deployed_at DESC
