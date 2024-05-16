DELETE FROM environments
WHERE id = @envID AND project_id = @projectID
