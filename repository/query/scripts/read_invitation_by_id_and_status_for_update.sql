SELECT *
FROM project_invitations
WHERE
    id = @id AND
    status = @status
FOR UPDATE
