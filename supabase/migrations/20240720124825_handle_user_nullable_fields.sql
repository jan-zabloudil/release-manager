-- If the user is created via SSO registration, such as Google SSO, the user name and avatar URL are provided.
-- For other registration methods, the user name and avatar URL may not be provided.
-- This migration updates the handle_new_user() and handle_user_update functions to insert an empty string for name and avatar URL instead of NULL.
-- It also updates the public.users table to set the name and avatar_url columns to NOT NULL.
-- This change ensures we don't have to worry about NULL values when scanning user data into Go structs.

CREATE OR REPLACE FUNCTION public.handle_new_user()
    RETURNS trigger
    LANGUAGE plpgsql
    SECURITY DEFINER
    SET search_path TO 'public'
AS $function$
BEGIN
    INSERT INTO public.users (id, email, name, avatar_url, created_at, updated_at)
    VALUES (
               new.id,
               new.email,
               COALESCE(new.raw_user_meta_data ->> 'name', ''),
               COALESCE(new.raw_user_meta_data ->> 'picture', ''),
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

    RETURN new;
END;
$function$;

CREATE OR REPLACE FUNCTION public.handle_user_update()
    RETURNS trigger
    LANGUAGE plpgsql
    SECURITY DEFINER
    SET search_path TO 'public'
AS $function$
BEGIN
    UPDATE public.users
    SET email = new.email,
        name = COALESCE(new.raw_user_meta_data ->> 'name', ''),
        avatar_url = COALESCE(new.raw_user_meta_data ->> 'picture', ''),
        updated_at = new.updated_at
    WHERE id = new.id;
    RETURN new;
END;
$function$;

UPDATE public.users
SET
    name = COALESCE(name, ''),
    avatar_url = COALESCE(avatar_url, '')
WHERE name IS NULL OR avatar_url IS NULL;

ALTER TABLE public.users
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN avatar_url SET NOT NULL;
