ALTER TABLE public.projects
ADD COLUMN github_owner_slug TEXT,
ADD COLUMN github_repo_slug TEXT,
ADD COLUMN github_repo_url TEXT;
