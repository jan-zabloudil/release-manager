SELECT *
FROM environments
WHERE project_id = @projectID
ORDER BY name
