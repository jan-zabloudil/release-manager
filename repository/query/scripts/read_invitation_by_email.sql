SELECT *
FROM project_invitations
WHERE
    project_id = @projectID AND
    email = @email
