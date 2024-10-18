SELECT *
FROM project_invitations
WHERE
    token_hash = @hash
