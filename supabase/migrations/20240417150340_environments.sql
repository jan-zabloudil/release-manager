CREATE TABLE public.environments (
         id UUID NOT NULL PRIMARY KEY,
         project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
         name TEXT NOT NULL,
         service_url TEXT,
         created_at TIMESTAMP WITH TIME ZONE NOT NULL,
         updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
         CONSTRAINT unique_environment_name_per_project UNIQUE (project_id, name)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE public.projects TO service_role;
