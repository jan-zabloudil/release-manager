DELETE FROM releases AS r
USING projects AS p
WHERE
    p.id = r.project_id AND
    p.github_owner_slug = @ownerSlug AND
    p.github_repo_slug = @repoSlug AND
    r.git_tag_name = @gitTagName
