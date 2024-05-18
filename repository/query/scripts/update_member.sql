UPDATE project_members
SET
    project_role = @projectRole,
    updated_at = @updatedAt
WHERE
    project_id = @projectID AND
    user_id = @userID
