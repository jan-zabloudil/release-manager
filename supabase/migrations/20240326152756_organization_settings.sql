create table "public"."organization_settings" (
                                                  "key" text not null,
                                                  "value" text not null
);


alter table "public"."organization_settings" enable row level security;

CREATE UNIQUE INDEX organization_settings_key_key ON public.organization_settings USING btree (key);

CREATE UNIQUE INDEX organization_settings_pkey ON public.organization_settings USING btree (key);

alter table "public"."organization_settings" add constraint "organization_settings_pkey" PRIMARY KEY using index "organization_settings_pkey";

alter table "public"."organization_settings" add constraint "organization_settings_key_key" UNIQUE using index "organization_settings_key_key";

grant delete on table "public"."organization_settings" to "service_role";

grant insert on table "public"."organization_settings" to "service_role";

grant references on table "public"."organization_settings" to "service_role";

grant select on table "public"."organization_settings" to "service_role";

grant trigger on table "public"."organization_settings" to "service_role";

grant truncate on table "public"."organization_settings" to "service_role";

grant update on table "public"."organization_settings" to "service_role";
