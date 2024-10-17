-- after refactoring of accept invitation logic in the service layer, this trigger and function are no longer needed

DROP TRIGGER IF EXISTS on_project_invitations_status_update ON public.project_invitations;
DROP FUNCTION IF EXISTS check_accepted_invitations_for_registered_users();
