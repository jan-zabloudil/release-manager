create table "public"."apps" (
                                 "id" uuid not null default gen_random_uuid(),
                                 "project_id" uuid not null,
                                 "name" text not null,
                                 "description" text,
                                 "environments" json,
                                 "created_at" timestamp with time zone not null default now(),
                                 "updated_at" timestamp with time zone not null default now()
);


CREATE UNIQUE INDEX apps_pkey ON public.apps USING btree (id);

alter table "public"."apps" add constraint "apps_pkey" PRIMARY KEY using index "apps_pkey";

alter table "public"."apps" add constraint "apps_project_id_fkey" FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE not valid;

alter table "public"."apps" validate constraint "apps_project_id_fkey";

grant delete on table "public"."apps" to "service_role";

grant insert on table "public"."apps" to "service_role";

grant references on table "public"."apps" to "service_role";

grant select on table "public"."apps" to "service_role";

grant trigger on table "public"."apps" to "service_role";

grant truncate on table "public"."apps" to "service_role";

grant update on table "public"."apps" to "service_role";
