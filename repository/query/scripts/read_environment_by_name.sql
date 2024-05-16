SELECT *
FROM environments
WHERE project_id = @projectID AND name = @name
