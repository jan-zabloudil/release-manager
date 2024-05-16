DELETE FROM releases
WHERE
	project_id = @projectID AND
	id = @releaseID
