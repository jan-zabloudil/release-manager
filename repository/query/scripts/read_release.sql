SELECT *
FROM releases
WHERE
    id = @releaseID AND
    project_id = @projectID
