UPDATE environments
SET
    name = @name,
    service_url = @serviceURL,
    updated_at = @updatedAt
WHERE
    id = @envID
