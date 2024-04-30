CREATE TABLE public.project_members (
    user_id UUID NOT NULL REFERENCES public.users ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
    project_role project_role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    PRIMARY KEY (user_id, project_id)
)
