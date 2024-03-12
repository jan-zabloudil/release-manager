create table "public"."projects_invitations" (
                                                 "id" uuid not null default gen_random_uuid(),
                                                 "project_id" uuid not null,
                                                 "email" text not null,
                                                 "role" text not null,
                                                 "invited_by_user_id" uuid,
                                                 "created_at" timestamp without time zone not null default now()
);
create table "public"."projects_members" (
                                             "project_id" uuid not null,
                                             "user_id" uuid not null,
                                             "role" text not null,
                                             "invited_by_user_id" uuid,
                                             "created_at" timestamp without time zone not null default now(),
                                             "updated_at" timestamp without time zone not null default now()
);


alter table "public"."projects_members" enable row level security;
alter table "public"."projects_invitations" enable row level security;

create policy "projects_invitations"
on "public"."projects_invitations"
as permissive
for all
to service_role;

create policy "projects_members_policy"
on "public"."projects_members"
as permissive
for all
to service_role;

CREATE UNIQUE INDEX projects_invitations_pkey ON public.projects_invitations USING btree (id);
CREATE UNIQUE INDEX projects_members_pkey ON public.projects_members USING btree (project_id, user_id);
CREATE UNIQUE INDEX unique_invitation ON public.projects_invitations USING btree (project_id, email);
CREATE UNIQUE INDEX unique_member ON public.projects_members USING btree (project_id, user_id);

alter table "public"."projects_invitations" add constraint "projects_invitations_pkey" PRIMARY KEY using index "projects_invitations_pkey";
alter table "public"."projects_members" add constraint "projects_members_pkey" PRIMARY KEY using index "projects_members_pkey";
alter table "public"."projects_invitations" add constraint "projects_invitations_invited_by_user_id_fkey" FOREIGN KEY (invited_by_user_id) REFERENCES auth.users(id) ON DELETE SET NULL not valid;
alter table "public"."projects_invitations" validate constraint "projects_invitations_invited_by_user_id_fkey";
alter table "public"."projects_invitations" add constraint "projects_invitations_project_id_fkey" FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE not valid;
alter table "public"."projects_invitations" validate constraint "projects_invitations_project_id_fkey";
alter table "public"."projects_invitations" add constraint "unique_invitation" UNIQUE using index "unique_invitation";
alter table "public"."projects_members" add constraint "projects_members_invited_by_user_id_fkey" FOREIGN KEY (invited_by_user_id) REFERENCES auth.users(id) ON DELETE SET NULL not valid;
alter table "public"."projects_members" validate constraint "projects_members_invited_by_user_id_fkey";
alter table "public"."projects_members" add constraint "projects_members_project_id_fkey" FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE not valid;
alter table "public"."projects_members" validate constraint "projects_members_project_id_fkey";
alter table "public"."projects_members" add constraint "projects_members_user_id_fkey" FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE not valid;
alter table "public"."projects_members" validate constraint "projects_members_user_id_fkey";
alter table "public"."projects_members" add constraint "unique_member" UNIQUE using index "unique_member";


CREATE OR REPLACE FUNCTION public.create_project_and_assign_member(name text, description text, slack_channel_id text, user_id uuid, role text)
 RETURNS SETOF projects
 LANGUAGE plpgsql
AS $function$
DECLARE
project_id UUID;
BEGIN
INSERT INTO projects (name, description, slack_channel_id)
VALUES (name, description, slack_channel_id)
    RETURNING id INTO project_id;

INSERT INTO projects_members (project_id, user_id, role)
VALUES (project_id, user_id, role);

RETURN QUERY SELECT * FROM projects WHERE id = project_id;
END;
$function$
;

CREATE OR REPLACE FUNCTION public.handle_new_user()
 RETURNS trigger
 LANGUAGE plpgsql
 SECURITY DEFINER
 SET search_path TO 'public'
AS $function$
begin
insert into public.projects_members (project_id, user_id, role, invited_by_user_id)
select pi.project_id, new.id, pi.role, pi.invited_by_user_id
from public.projects_invitations pi
where pi.email = new.email;

delete from public.projects_invitations
where email = new.email;

return new;
end;
$function$
;

create trigger on_auth_user_created
    after insert on auth.users
    for each row execute procedure public.handle_new_user();
CREATE TRIGGER projects_members_handle_updated_at BEFORE UPDATE ON public.projects_members FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');

CREATE OR REPLACE FUNCTION public.get_project_member(p_project_id uuid, p_user_id uuid)
 RETURNS TABLE(user_data jsonb, project_id uuid, role text, invited_by_user_id uuid, created_at timestamp without time zone, updated_at timestamp without time zone)
 LANGUAGE plpgsql
 STABLE SECURITY DEFINER
AS $function$
BEGIN
RETURN QUERY
SELECT
    jsonb_build_object(
            'id', u.id,
            'aud', u.aud,
            'role', u.role,
            'email', u.email,
            'invited_at', u.invited_at,
            'confirmed_at', u.confirmed_at,
            'confirmation_sent_at', u.confirmation_sent_at,
            'app_metadata', u.raw_app_meta_data,
            'user_metadata', u.raw_user_meta_data,
            'created_at', u.created_at,
            'updated_at', u.updated_at
    ) AS user_data,
    m.project_id,
    m.role,
    m.invited_by_user_id,
    m.created_at,
    m.updated_at
FROM auth.users u
         JOIN public.projects_members m ON m.user_id = u.id
WHERE m.project_id = p_project_id
  AND u.id = p_user_id;
END;
$function$
;

CREATE OR REPLACE FUNCTION public.get_project_members(p_project_id uuid)
 RETURNS TABLE(user_data jsonb, project_id uuid, role text, invited_by_user_id uuid, created_at timestamp without time zone, updated_at timestamp without time zone)
 LANGUAGE plpgsql
 STABLE SECURITY DEFINER
AS $function$
BEGIN
RETURN QUERY
SELECT
    jsonb_build_object(
            'id', u.id,
            'aud', u.aud,
            'role', u.role,
            'email', u.email,
            'invited_at', u.invited_at,
            'confirmed_at', u.confirmed_at,
            'confirmation_sent_at', u.confirmation_sent_at,
            'app_metadata', u.raw_app_meta_data,
            'user_metadata', u.raw_user_meta_data,
            'created_at', u.created_at,
            'updated_at', u.updated_at
    ) AS user_data,
    m.project_id,
    m.role,
    m.invited_by_user_id,
    m.created_at,
    m.updated_at
FROM auth.users u
         JOIN public.projects_members m ON m.user_id = u.id
WHERE m.project_id = p_project_id;
END;
$function$
;

CREATE OR REPLACE FUNCTION public.get_user_by_email(p_email text)
 RETURNS TABLE(id uuid, aud character varying, role character varying, email character varying, invited_at timestamp with time zone, confirmed_at timestamp with time zone, confirmation_sent_at timestamp with time zone, app_metadata jsonb, user_metadata jsonb, created_at timestamp with time zone, updated_at timestamp with time zone)
 LANGUAGE plpgsql
 STABLE SECURITY DEFINER
AS $function$
BEGIN
RETURN QUERY
SELECT
    u.id,
    u.aud,
    u.role,
    u.email,
    u.invited_at,
    u.confirmed_at,
    u.confirmation_sent_at,
    u.raw_app_meta_data AS app_metadata,
    u.raw_user_meta_data AS user_metadata,
    u.created_at,
    u.updated_at
FROM auth.users u
WHERE u.email = p_email;
END;
$function$
;

CREATE OR REPLACE FUNCTION public.create_project_invitation(
    p_project_id UUID,
    p_email TEXT,
    p_role TEXT,
    p_invited_by_user_id UUID
)
RETURNS SETOF public.projects_invitations
LANGUAGE plpgsql
VOLATILE SECURITY DEFINER
AS $function$
BEGIN
    IF EXISTS (
        SELECT 1 FROM auth.users WHERE email = p_email
    ) THEN
        RAISE EXCEPTION 'A user with the provided email (%) exists. Invitations should only be sent to users who have not signed up yet.', p_email USING ERRCODE = 'P0001';
END IF;

INSERT INTO public.projects_invitations(project_id, email, role, invited_by_user_id)
VALUES (p_project_id, p_email, p_role, p_invited_by_user_id);

RETURN QUERY SELECT * FROM public.projects_invitations
    WHERE email = p_email AND project_id = p_project_id;

END;
$function$;

