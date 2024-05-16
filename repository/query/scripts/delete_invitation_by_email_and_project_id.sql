DELETE FROM project_invitations
WHERE email = @email AND project_id = @projectID
