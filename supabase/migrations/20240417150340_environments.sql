create table "public"."environments" (
         "id" uuid not null default gen_random_uuid(),
         "project_id" uuid not null references public.projects on delete cascade,
         "name" text not null,
         "service_url" text,
         "created_at" timestamp with time zone not null default now(),
         "updated_at" timestamp with time zone not null default now(),
         CONSTRAINT environments_pkey PRIMARY KEY (id),
         CONSTRAINT unique_name_per_project UNIQUE (project_id, name)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE public.projects TO service_role;

