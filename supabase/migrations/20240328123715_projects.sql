CREATE TABLE public.projects (
     id uuid NOT NULL DEFAULT gen_random_uuid(),
     name text NOT NULL,
     slack_channel_id text,
     release_notification_config json,
     created_at timestamp with time zone NOT NULL,
     updated_at timestamp with time zone NOT NULL,
     CONSTRAINT projects_pkey PRIMARY KEY (id)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE public.projects TO service_role;
