-- Although the application logic does not require following fields to always be set, if they are missing in the input,
-- they will be inserted as empty strings.
-- And when reading from the database, the application do not expect these fields to be NULL.
-- This migration ensures that the expected behavior is reflected in the database schema.

UPDATE public.projects
SET slack_channel_id = ''
WHERE slack_channel_id IS NULL;

UPDATE public.projects
SET release_notification_config = '{}'::json
WHERE release_notification_config IS NULL;

ALTER TABLE public.projects
ALTER COLUMN slack_channel_id SET NOT NULL,
ALTER COLUMN release_notification_config SET NOT NULL;

UPDATE public.environments
SET service_url = ''
WHERE service_url IS NULL;

ALTER TABLE public.environments
ALTER COLUMN service_url SET NOT NULL;

UPDATE public.releases
SET release_notes = ''
WHERE release_notes IS NULL;

ALTER TABLE public.releases
ALTER COLUMN release_notes SET NOT NULL;
