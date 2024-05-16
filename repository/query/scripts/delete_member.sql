DELETE FROM project_members
WHERE project_id = @projectID AND user_id = @userID
