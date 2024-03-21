alter table "public"."apps" add column "scm_repo" json;
CREATE TRIGGER apps_handle_updated_at BEFORE UPDATE ON public.apps FOR EACH ROW EXECUTE FUNCTION moddatetime('updated_at');
