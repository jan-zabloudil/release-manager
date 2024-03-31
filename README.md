# ReleaseManager

## Running locally without Docker

### Supabase

Supabase is used for authentication with Google OAuth and storing data in PostgreSQL database. 

1. Install [Supabase CLI](https://supabase.com/docs/guides/cli/getting-started)
2. To configure `Google OAuth` using Supabase, set the following environment variables _(if `.env` file is present, env variables are set from the file)_:

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

| Env variable          | Where do I find this value?                                                                                                                                                           | Default |
|-----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `PORT`                | -                                                                                                                                                                                     | `8080`  |
| `LOG_LEVEL`           | Allowed values `DEBUG`, `INFO`, `WARN`, `ERROR`.                                                                                                                                      | `INFO`  |
| `SUPABASE_API_URL`    | Assuming that the local Supabase instance is running, navigate to the project folder and run `supabase status` in your terminal. The Supabase credentials are printed, use `API URL`. | -       |
| `SUPABASE_SECRET_KEY` | As mentioned above, run `supabase status` and use `service_role key`.                                                                                                                 | -       |

> If you are using hosted Supabase, navigate to Supabase Studio, then go to *Your project > Project Settings > API* to find the api url and secret key.

2. Run app `make run_local`

## How to create admin user?

1. Sign up as a regular user.
2. Navigate to Supabase Studio.
3. Run the following query to add the `is_admin` flag to the user.
```
select set_admin_user('user_id');
```
_User ID can be found in `auth.users`._

