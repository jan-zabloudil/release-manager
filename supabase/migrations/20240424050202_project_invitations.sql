BEGIN;

CREATE TYPE project_role AS ENUM ('editor', 'viewer');
CREATE TYPE project_invitation_status AS ENUM ('pending', 'accepted_awaiting_registration');

CREATE TABLE public.project_invitations (
    id UUID NOT NULL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
    email TEXT NOT NULL,
    project_role project_role NOT NULL,
    status project_invitation_status NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    invited_by UUID REFERENCES public.users ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT unique_invitation_per_project UNIQUE (project_id, email)
);

COMMIT;

