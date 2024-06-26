UPDATE projects
SET
    name = @name,
    slack_channel_id = @slackChannelID,
    release_notification_config = @releaseNotificationConfig,
    github_owner_slug = @githubOwnerSlug,
    github_repo_slug = @githubRepoSlug,
    updated_at = @updatedAt
WHERE id = @id
