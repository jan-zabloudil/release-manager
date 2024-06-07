ALTER TABLE public.projects
ADD COLUMN github_owner_slug TEXT,
ADD COLUMN github_repo_slug TEXT,
ADD CONSTRAINT unique_github_repo UNIQUE (github_owner_slug, github_repo_slug)
