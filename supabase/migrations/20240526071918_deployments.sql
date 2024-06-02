CREATE TABLE public.deployments (
     id UUID NOT NULL PRIMARY KEY,
     release_id UUID NOT NULL REFERENCES public.releases ON DELETE CASCADE,
     environment_id UUID NOT NULL REFERENCES public.environments ON DELETE CASCADE,
     deployed_by UUID REFERENCES public.users ON DELETE SET NULL,
     deployed_at TIMESTAMP WITH TIME ZONE NOT NULL
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
    ON TABLE public.deployments TO service_role;
