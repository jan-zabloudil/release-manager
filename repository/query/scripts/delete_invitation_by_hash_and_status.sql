DELETE FROM project_invitations
WHERE token_hash = @hash AND status = @status
