CREATE TABLE public.releases (
     id UUID NOT NULL PRIMARY KEY,
     project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
     release_title TEXT NOT NULL,
     release_notes TEXT,
     created_by UUID REFERENCES public.users ON DELETE SET NULL,
     created_at TIMESTAMP WITH TIME ZONE NOT NULL,
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
     CONSTRAINT unique_release_title_per_project UNIQUE (project_id, release_title)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
    ON TABLE public.releases TO service_role;
