create table "public"."templates" (
      "id" uuid not null default gen_random_uuid(),
      "type" text not null,
      "template_data" json not null
);

CREATE UNIQUE INDEX templates_pkey ON public.templates USING btree (id);

alter table "public"."templates" add constraint "templates_pkey" PRIMARY KEY using index "templates_pkey";

grant delete on table "public"."templates" to "service_role";
grant insert on table "public"."templates" to "service_role";
grant references on table "public"."templates" to "service_role";
grant select on table "public"."templates" to "service_role";
grant trigger on table "public"."templates" to "service_role";
grant truncate on table "public"."templates" to "service_role";
grant update on table "public"."templates" to "service_role";
