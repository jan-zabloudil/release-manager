CREATE TABLE public.release_attachments (
    release_id UUID NOT NULL REFERENCES public.releases ON DELETE CASCADE,
    attachment_id UUID NOT NULL,
    name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (release_id, attachment_id),
    CONSTRAINT unique_file_path UNIQUE (file_path)
)
