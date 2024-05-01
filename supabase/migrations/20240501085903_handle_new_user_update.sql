-- extends function added in 20240331171034_public_users.sql
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
               new.raw_user_meta_data ->> 'name',
               new.raw_user_meta_data ->> 'picture',
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
