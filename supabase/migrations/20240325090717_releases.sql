create table "public"."releases" (
                                     "app_id" uuid not null,
                                     "source_code" json not null,
                                     "deployments" json,
                                     "changelog" text,
                                     "created_by_user_id" uuid not null,
                                     "created_at" timestamp with time zone not null default now(),
                                     "updated_at" timestamp with time zone not null default now(),
                                     "id" uuid not null default gen_random_uuid(),
                                     "title" text not null
);


CREATE UNIQUE INDEX releases_pkey ON public.releases USING btree (id);

alter table "public"."releases" add constraint "releases_pkey" PRIMARY KEY using index "releases_pkey";

alter table "public"."releases" add constraint "releases_app_id_fkey" FOREIGN KEY (app_id) REFERENCES apps(id) not valid;

alter table "public"."releases" validate constraint "releases_app_id_fkey";

alter table "public"."releases" add constraint "releases_created_by_user_id_fkey" FOREIGN KEY (created_by_user_id) REFERENCES auth.users(id) not valid;

alter table "public"."releases" validate constraint "releases_created_by_user_id_fkey";

grant delete on table "public"."releases" to "service_role";

grant insert on table "public"."releases" to "service_role";

grant references on table "public"."releases" to "service_role";

grant select on table "public"."releases" to "service_role";

grant trigger on table "public"."releases" to "service_role";

grant truncate on table "public"."releases" to "service_role";

grant update on table "public"."releases" to "service_role";
