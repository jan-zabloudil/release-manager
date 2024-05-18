UPDATE project_invitations
SET
    status = @status,
    updated_at = @updatedAt
WHERE
    id = @invitationID
