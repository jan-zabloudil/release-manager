CREATE TABLE public.releases (
     id UUID NOT NULL PRIMARY KEY,
     project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
     release_title TEXT NOT NULL,
     release_notes TEXT,
     created_by UUID REFERENCES public.users ON DELETE SET NULL,
     created_at TIMESTAMP WITH TIME ZONE NOT NULL,
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
     git_tag_name Text NOT NULL,
     github_release_created_at TIMESTAMP WITH TIME ZONE NOT NULL,
     github_release_published_at TIMESTAMP WITH TIME ZONE NOT NULL,
     CONSTRAINT unique_git_tag_per_project UNIQUE (project_id, git_tag_name)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
    ON TABLE public.releases TO service_role;
