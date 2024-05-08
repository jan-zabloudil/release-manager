CREATE OR REPLACE FUNCTION create_project_and_owner_member (
    p_id UUID,
    p_name TEXT,
    p_slack_channel_id TEXT,
    p_release_notification_config JSON,
    p_github_repository JSON,
    p_user_id UUID,
    p_project_role project_role,
    p_created_at TIMESTAMP WITH TIME ZONE,
    p_updated_at TIMESTAMP WITH TIME ZONE
)
    RETURNS VOID AS $$
BEGIN
    INSERT INTO public.projects (id, name, slack_channel_id, release_notification_config, github_repository, created_at, updated_at)
    VALUES (p_id, p_name, p_slack_channel_id, p_release_notification_config, p_github_repository, p_created_at, p_updated_at);

    INSERT INTO public.project_members (user_id, project_id, project_role, created_at, updated_at)
    VALUES (p_user_id, p_id, p_project_role, p_created_at, p_updated_at);
END;
$$ LANGUAGE plpgsql;
