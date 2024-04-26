CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE public.users (
      id UUID NOT NULL PRIMARY KEY REFERENCES auth.users ON DELETE CASCADE,
      email TEXT NOT NULL UNIQUE,
      name TEXT,
      avatar_url TEXT,
      role user_role NOT NULL DEFAULT 'user',
      created_at TIMESTAMP WITH TIME ZONE NOT NULL,
      updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE "public"."users" TO "service_role";

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
    RETURN new;
END;
$function$;

CREATE TRIGGER on_auth_user_created
    AFTER INSERT ON auth.users
    FOR EACH ROW EXECUTE PROCEDURE public.handle_new_user();

CREATE OR REPLACE FUNCTION public.handle_user_update()
    RETURNS trigger
    LANGUAGE plpgsql
    SECURITY DEFINER
    SET search_path TO 'public'
AS $function$
BEGIN
    UPDATE public.users
    SET email = new.email,
        name = new.raw_user_meta_data ->> 'name',
        avatar_url = new.raw_user_meta_data ->> 'picture',
        updated_at = new.updated_at
    WHERE id = new.id;
    RETURN new;
END;
$function$;

CREATE TRIGGER on_auth_user_updated
    AFTER UPDATE ON auth.users
    FOR EACH ROW EXECUTE PROCEDURE public.handle_user_update();
