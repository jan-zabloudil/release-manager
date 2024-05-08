CREATE OR REPLACE FUNCTION create_project_member_and_delete_invitation(
    p_user_id UUID,
    p_project_id UUID,
    p_email TEXT,
    p_project_role project_role,
    p_created_at TIMESTAMP WITH TIME ZONE,
    p_updated_at TIMESTAMP WITH TIME ZONE
)
    RETURNS void AS $$
BEGIN
        INSERT INTO public.project_members (user_id, project_id, project_role, created_at, updated_at)
        VALUES (p_user_id, p_project_id, p_project_role, p_created_at, p_updated_at);

        DELETE FROM public.project_invitations
        WHERE project_id = p_project_id AND email = p_email;
END;
$$ LANGUAGE plpgsql;
