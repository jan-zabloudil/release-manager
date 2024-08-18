-- Migration creates a table to store pending admin user role assignments
-- and adds a function (and trigger) that will prevent of adding email to the pending_admin_user_role_assignments table if the user with this email already exists.

-- Migration also updates the handle_new_user function. The function now assigns the user role based on the presence of user's email in the pending_admin_user_role_assignments table.
-- IMPORTANT: The function do not check if the user is verified.
-- If login via email is allowed in the future, make sure you allow login only to verified users (can be set in Supabase),
-- otherwise anyone could sing up with an email that is in the pending_admin_user_role_assignments table and get the admin role.

-- Motivation:
-- Before this migration, all users were created with the 'user' role.
-- To set a user as an admin, the user had to sing up first and then an superadmin of the platform (with access to the database) could manually update the user role.
-- That was very annoying and slow process.
-- With this migration, a superadmin can set all admins user before they sign up. And once they sign up, they can use the platform with the admin role straight away.

CREATE TABLE public.pending_admin_user_role_assignments (
    email TEXT PRIMARY KEY
);

CREATE OR REPLACE FUNCTION public.handle_new_user()
    RETURNS trigger
    LANGUAGE plpgsql
    SECURITY DEFINER
    SET search_path TO 'public'
AS $function$
DECLARE
    role_to_assign user_role := 'user';
BEGIN
    IF EXISTS (
        SELECT 1
        FROM public.pending_admin_user_role_assignments
        WHERE
            email = NEW.email AND
          -- Added additional security check to prevent assigning the admin role to unverified users (If login via email is allowed in the future).
          -- Column raw_user_meta_data is used because auth.users.email_confirmed_at is NULL at point of inserting into the auth.users table.

          -- IMPORTANT: when inserting users manually via Supabase UI, raw_user_meta_data is not set.
          -- In such case, set admin role manually.
            NEW.raw_user_meta_data ->> 'email_verified' = 'true'
    ) THEN
        role_to_assign := 'admin';
    END IF;

    INSERT INTO public.users (id, email, name, avatar_url, role, created_at, updated_at)
    VALUES (
               new.id,
               new.email,
               COALESCE(new.raw_user_meta_data ->> 'name', ''),
               COALESCE(new.raw_user_meta_data ->> 'picture', ''),
               role_to_assign,
               new.created_at,
               new.updated_at
           );

    -- If user has accepted an invitation(s) before signing up, add them to the project(s)
    INSERT INTO public.project_members (user_id, project_id, project_role, created_at, updated_at)
    SELECT NEW.id, project_id, project_role, NEW.created_at, NEW.updated_at
    FROM public.project_invitations
    WHERE email = NEW.email AND status = 'accepted_awaiting_registration';

    DELETE FROM public.project_invitations
    WHERE email = NEW.email AND status = 'accepted_awaiting_registration';

    DELETE FROM public.pending_admin_user_role_assignments
    WHERE email = NEW.email;

    RETURN new;
END;
$function$;

-- Add function (and trigger) that will prevent of adding email to the pending_admin_user_role_assignments table if the user with this email already exists.
CREATE OR REPLACE FUNCTION public.check_existing_user()
    RETURNS trigger
    LANGUAGE plpgsql
AS $function$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM public.users
        WHERE email = NEW.email
    ) THEN
        RAISE EXCEPTION 'Email % already exists in the users table. Cannot insert into pending_admin_user_role_assignments.', NEW.email;
    END IF;

    RETURN NEW;
END;
$function$;

CREATE TRIGGER check_existing_user
    BEFORE INSERT ON public.pending_admin_user_role_assignments
    FOR EACH ROW
EXECUTE FUNCTION public.check_existing_user();
