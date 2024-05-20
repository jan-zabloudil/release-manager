INSERT INTO project_invitations (id, project_id, email, project_role, status, token_hash, invited_by, created_at, updated_at)
VALUES (@invitationID, @projectID, @email, @projectRole, @status, @tokenHash, @invitedBy, @createdAt, @updatedAt)
