SELECT
    u.id,
    u.email,
    u.name,
    u.avatar_url,
    u.role,
    u.created_at,
    u.updated_at,
    pm.project_id,
    pm.project_role,
    pm.created_at,
    pm.updated_at
FROM project_members pm
JOIN users u ON u.id = pm.user_id
WHERE
    pm.project_id = @projectID
ORDER BY u.name
