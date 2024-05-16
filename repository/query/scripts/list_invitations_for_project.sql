SELECT *
FROM project_invitations
WHERE project_id = @projectID
ORDER BY created_at DESC
