SELECT *
FROM releases
WHERE
    release_title = @releaseTitle AND
    project_id = @projectID
