CREATE TABLE public.projects (
     id UUID NOT NULL PRIMARY KEY,
     name TEXT NOT NULL,
     slack_channel_id TEXT,
     release_notification_config JSON,
     created_at TIMESTAMP WITH TIME ZONE NOT NULL,
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE public.projects TO service_role;
