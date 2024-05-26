CREATE TABLE public.releases (
     id UUID NOT NULL PRIMARY KEY,
     project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
     release_title TEXT NOT NULL,
     release_notes TEXT,
     created_by UUID REFERENCES public.users ON DELETE SET NULL,
     created_at TIMESTAMP WITH TIME ZONE NOT NULL,
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
     github_release_id BIGINT NOT NULL,
     github_owner_slug TEXT NOT NULL,
     github_repo_slug TEXT NOT NULL,
     github_release_data JSON NOT NULL,
     CONSTRAINT unique_github_release_per_project UNIQUE (project_id, github_owner_slug, github_repo_slug, github_release_id)
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
    ON TABLE public.releases TO service_role;
