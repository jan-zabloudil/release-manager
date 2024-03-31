create table "public"."users" (
      id uuid not null references auth.users on delete cascade,
      email text not null UNIQUE,
      name text,
      avatar_url text,
      role text not null DEFAULT 'user',
      created_at timestamp with time zone,
      updated_at timestamp with time zone,

      primary key (id)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE "public"."users" TO "service_role";

CREATE OR REPLACE FUNCTION set_admin_user(user_id uuid)
RETURNS void AS
$$
BEGIN
    PERFORM set_claim(user_id, 'role', '"admin"');

    UPDATE public.users
    SET role = 'admin'
    WHERE id = user_id;

    IF NOT FOUND THEN
            RAISE EXCEPTION 'User with ID % not found.', user_id;
    END IF;

    EXCEPTION
        WHEN OTHERS THEN
            RAISE;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.handle_new_user()
 RETURNS trigger
 LANGUAGE plpgsql
 SECURITY DEFINER
 SET search_path TO 'public'
AS $function$
begin
insert into public.users (id, email, name, avatar_url, created_at, updated_at)
values (
           new.id,
           new.email,
           new.raw_user_meta_data ->> 'name',
           new.raw_user_meta_data ->> 'picture',
           new.created_at,
           new.updated_at
       );
return new;
end;
$function$
;

create trigger on_auth_user_created
    after insert on auth.users
    for each row execute procedure public.handle_new_user();

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
