-- trigger to prevent updating the status of a project invitation to 'accepted_awaiting_registration' if a user with the same email already exists
-- this race condition can happen in ProjectService -> AcceptInvitation(), scenario:
-- 1. service checks if user exists, user does not exist yet
-- 2. service accepts invitation, status is updated to 'accepted_awaiting_registration'
-- but the first and second steps are not executed in the same transaction, so a user with the same email can be created in the meantime
-- very unlikely, but still possible

CREATE OR REPLACE FUNCTION check_accepted_invitations_for_registered_users()
    RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.status = 'accepted_awaiting_registration') AND EXISTS (
        SELECT 1
        FROM public.users AS u
        WHERE u.email = NEW.email
    ) THEN
        RAISE EXCEPTION 'Cannot update status to accepted_awaiting_registration because a user with email % already exists', NEW.email;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER on_project_invitations_status_update
    BEFORE UPDATE OF status ON public.project_invitations
    FOR EACH ROW EXECUTE FUNCTION check_accepted_invitations_for_registered_users();
