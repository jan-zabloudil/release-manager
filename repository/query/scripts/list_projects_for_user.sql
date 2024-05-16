SELECT p.*
FROM projects p
JOIN project_members pm
    ON p.id = pm.project_id
WHERE pm.user_id = @userID
