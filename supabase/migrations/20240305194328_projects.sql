create table "public"."projects" (
                                     "id" uuid not null default gen_random_uuid(),
                                     "name" text not null default ''::text,
                                     "description" text,
                                     "slack_channel_id" text,
                                     "release_message_template" text,
                                     "created_at" timestamp without time zone not null default now(),
                                     "updated_at" timestamp without time zone not null default now()
);

CREATE UNIQUE INDEX projects_pkey ON public.projects USING btree (id);
alter table "public"."projects" add constraint "projects_pkey" PRIMARY KEY using index "projects_pkey";

alter table "public"."projects" enable row level security;
CREATE POLICY "projects_policy" ON "public"."projects"
AS PERMISSIVE FOR ALL
TO service_role;

create extension if not exists "moddatetime" with schema "extensions";
CREATE TRIGGER projects_handle_updated_at BEFORE UPDATE ON public.projects FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');
