CREATE TABLE public.settings (
     key TEXT NOT NULL PRIMARY KEY,
     value JSON NOT NULL
);

GRANT DELETE, INSERT, REFERENCES, SELECT, TRIGGER, TRUNCATE, UPDATE
ON TABLE public.settings TO service_role;

-- Application expects that these settings are present, therefore default values are inserted.
-- Settings can be changed via the application.
INSERT INTO public.settings (key, value) VALUES
    ('organization_name', '"Your company"'),
    ('default_release_message', '"Hey everyone, we are thrilled to announce new release!"'),
    ('slack', '{"enabled": false, "token": ""}'),
    ('github', '{"enabled": false, "token": ""}');
