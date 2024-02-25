
# ReleaseManager

## Running locally without Docker

### Supabase

Supabase is used for authentication with Google OAuth and storing data in PostgreSQL database.

1. Install [Supabase CLI](https://supabase.com/docs/guides/cli/getting-started)
2. To configure `Google OAuth` using Supabase, set the following environment variables:

| Env variable                                 | Where do I find this value?                                                                                                              |
|----------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| `SUPABASE_AUTH_EXTERNAL_GOOGLE_CLIENT_ID`    | [GCC](https://console.cloud.google.com/) > APIs & Services > Credentials > _Choose your client_ > Additional information > Client ID     |
| `SUPABASE_AUTH_EXTERNAL_GOOGLE_SECRECT`      | [GCC](https://console.cloud.google.com/) > APIs & Services > Credentials > _Choose your client_ > Additional information > Client secret |
| `SUPABASE_AUTH_EXTERNAL_GOOGLE_REDIRECT_URI` | [GCC](https://console.cloud.google.com/) > APIs & Services > Credentials > _Choose your client_ > Authorized redirect URIs               |

> How to create a Google OAuth client? See [official docs](https://support.google.com/cloud/answer/6158849?hl=en) or watch [Supabase tutorial](https://youtu.be/_XM9ziOzWk4?si=22ZX02UgJtHEXVtZ&t=25).
>
>If you are using hosted Supabase, navigate to Supabase Studio, then go to *Your project > Authentication > Providers > Google*, and configure Google Auth there.
3. Navigate to the project folder and start the local Supabase by running `supabase start`.
4. To test if Google OAuth was set up correctly, open your browser and navigate to `<SUPABASE_API_URL>/auth/v1/authorize?provider=google`.

### REST API

1. Now set env variables for REST API:

| Env variable                          | Where do I find this value?                                                                                                                                                           | Default |
|---------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `RELEASE_MANAGER_PORT`                | -                                                                                                                                                                                     | `8080`  |
| `RELEASE_MANAGER_LOG_LEVEL`           | Allowed values `DEBUG`, `INFO`, `WARN`, `ERROR`.                                                                                                                                      | `DEBUG` |
| `RELEASE_MANAGER_SUPABASE_API_URL`    | Assuming that the local Supabase instance is running, navigate to the project folder and run `supabase status` in your terminal. The Supabase credentials are printed, use `API URL`. | -       |
| `RELEASE_MANAGER_SUPABASE_SECRET_KEY` | As mentioned above, run `supabase status` and use `service_role key`.                                                                                                                 | -       |

> If you are using hosted Supabase, navigate to Supabase Studio, then go to *Your project > Project Settings > API* to find the api url and secret key.
2. Run app `make run_local`

## How to create admin user?
1. First sign up as regular user. 
2. Admin flag can be added to the user via the Supabase SQL editor by executing [custom_claim](https://github.com/supabase-community/supabase-custom-claims): `select set_claim('user_id', 'is_admin', 'true')`.

If the `is_admin` role is removed for a user, all user's tokens should be invalidated.
