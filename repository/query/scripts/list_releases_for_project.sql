SELECT *
FROM releases
WHERE project_id = @projectID
ORDER BY created_at DESC
