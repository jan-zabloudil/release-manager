ALTER TABLE public.environments
ADD COLUMN deployed_release_id UUID REFERENCES public.releases ON DELETE SET NULL;
