UPDATE projects
SET
    name = @name,
    slack_channel_id = @slackChannelID,
    release_notification_config = @releaseNotificationConfig,
    github_repository_url = @githubRepositoryURL,
    updated_at = @updatedAt
WHERE id = @id
