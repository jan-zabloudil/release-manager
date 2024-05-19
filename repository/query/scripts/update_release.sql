UPDATE releases
SET
    release_title = @releaseTitle,
    release_notes = @releaseNotes,
    updated_at = @updatedAt
WHERE
    id = @releaseID
