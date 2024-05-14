SELECT *
FROM projects
WHERE id = @id
FOR UPDATE
