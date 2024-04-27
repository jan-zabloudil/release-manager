BEGIN;

CREATE TYPE project_role AS ENUM ('owner', 'editor', 'viewer');
CREATE TYPE project_invitation_status AS ENUM ('pending', 'accepted_awaiting_registration');

-- Owner is a special role that is not assigned via invitations
CREATE DOMAIN invitation_project_role AS project_role
    CHECK (VALUE IN ('editor', 'viewer'));

CREATE TABLE public.project_invitations (
    id UUID NOT NULL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES public.projects ON DELETE CASCADE,
    email TEXT NOT NULL,
    project_role invitation_project_role NOT NULL,
    status project_invitation_status NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    invited_by UUID REFERENCES public.users ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT unique_invitation_per_project UNIQUE (project_id, email)
);

COMMIT;

