DELETE FROM release_attachments
WHERE
    release_id = @releaseID AND
    project_id = @projectID AND
    attachment_id = @attachmentID
