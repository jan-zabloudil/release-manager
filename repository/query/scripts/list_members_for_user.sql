SELECT
    u.id              AS user_id,
    u.email           AS user_email,
    u.name            AS user_name,
    u.avatar_url      AS user_avatar_url,
    u.role            AS user_role,
    u.created_at      AS user_created_at,
    u.updated_at      AS user_updated_at,
    pm.project_id     AS project_id,
    pm.project_role   AS project_role,
    pm.created_at     AS member_created_at,
    pm.updated_at     AS member_updated_at
FROM project_members pm
JOIN users u ON u.id = pm.user_id
WHERE
    u.id = @userID
ORDER BY u.name
