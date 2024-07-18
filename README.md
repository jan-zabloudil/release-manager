<div align="center">
  <h1>ReleaseManager</h1>
  <p>Platform-agnostic tool that streamlines daily operations in release management, planning and deployments.</p>
  <br />
  <img src="https://i.ibb.co/gF7QHDp/release-manager.png" alt="ReleaseManager" width="500">
</div>

## About

The vision is to develop a cloud and platform-agnostic tool that streamlines daily operations in release management, planning and deployments.

The current state of the project is a minimal viable product (MVP). More features to implement:

- [ ] Integration with AWS and GCP for deployment automation
- [ ] Release planning
- [ ] Integration with Jira for issue tracking

## Motivation

Currently, there aren’t many tools available, especially tools that are cloud and platform-agnostic, meaning that can be used to manage releases and deployments to various cloud providers, such as AWS or GCP, and that would be suitable for platforms like the backend, web or mobile.

## How to run the app?

To run the app, you need to set up the following services:

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

#### 1. Set env variables for REST API.

| Env variable                             | Where do I find this value?                                                                                                                                                                                                                                                                                                                                                                                      | Default |
|------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| `PORT`                                   | -                                                                                                                                                                                                                                                                                                                                                                                                                | `8080`  |
| `LOG_LEVEL`                              | Allowed values `DEBUG`, `INFO`, `WARN`, `ERROR`.                                                                                                                                                                                                                                                                                                                                                                 | `INFO`  |
| `SUPABASE_API_URL`                       | Assuming that the local Supabase instance is running, navigate to the project folder and run `supabase status` in your terminal. The Supabase credentials are printed, use `API URL`.                                                                                                                                                                                                                            | -       |
| `SUPABASE_API_SECRET_KEY`                | As mentioned above, run `supabase status` and use `service_role key`.                                                                                                                                                                                                                                                                                                                                            | -       |
| `SUPABASE_DATABASE_URL`                  | See [connection pooler](https://supabase.com/docs/guides/database/connecting-to-postgres#connection-pooler) or [direct connections](https://supabase.com/docs/guides/database/connecting-to-postgres#direct-connections) for details on obtaining a connection string.<br />If you are running Supabase locally, use the following connection string: `postgresql://postgres:postgres@localhost:54322/postgres`. | -       |
| `RESEND_API_KEY`                         | [Resend](https://resend.com/) is used for sending emails. Sign up to get API key. Login to your Resend account and create API key.                                                                                                                                                                                                                                                                               | -       |
| `RESEND_TEST_RECIPIENT`                  | Your Resend account email or there are several `*@resend.dev` options, see [docs](https://resend.com/docs/dashboard/emails/send-test-emails).                                                                                                                                                                                                                                                                    | -       |
| `RESEND_SENDER`                          | Sender email address. You have to verify your domain in order to be able to provide your email, for testing use `onboarding@resend.dev`                                                                                                                                                                                                                                                                          | -       |
| `RESEND_SEND_TO_REAL_RECIPIENTS`         | The value determines whether you want to send emails to real users or forward all emails to the Resend test recipient.                                                                                                                                                                                                                                                                                           | -       |
| `CLIENT_SERVICE_URL`                     | The URL where the client app is running.                                                                                                                                                                                                                                                                                                                                                                         | -       |
| `CLIENT_SERVICE_SIGN_UP_ROUTE`           | The route where the client app sign-up page is located.                                                                                                                                                                                                                                                                                                                                                          | -       |
| `CLIENT_SERVICE_ACCEPT_INVITATION_ROUTE` | The route where the client app accept invitation page is located.                                                                                                                                                                                                                                                                                                                                                | -       |
| `CLIENT_SERVICE_REJECT_INVITATION_ROUTE` | The route where the client app reject invitation page is located.                                                                                                                                                                                                                                                                                                                                                | -       |


> If you are using hosted Supabase, navigate to Supabase Studio, then go to *Your project > Project Settings > API* to find the api url and secret key. 

#### 2. Running the REST API

If you want to run the REST API locally, navigate to the project folder and run `make run_local`.

> To be able to run the REST API in Docker, you need to set `HOST_PORT` env variable.
>
> Host port is the port where the REST API will be available on your machine.

If you want to run the REST API in Docker, navigate to the project folder and run `make run_docker`.

If you want to rebuild the image and run the REST API, run `make run_docker_rebuild`.

## REST API Documentation

See [API documentation](api-doc.yaml).

## App configuration

Instructions on how to configure the app.

### How to create admin user?

1. Sign up as a regular user.
2. Navigate to Supabase Studio.
3. Go to *Your project > SQL Editor*.
3. Change `public.users.role` to `admin`.

### How to enable GitHub integration?

To enable GitHub integration, you need to call the REST API endpoint `PATCH /organization/settings` with the following payload:

```json
{
  "github": {
    "enabled": true,
    "token": "<GITHUB_TOKEN>",
    "webhook_secret": "<WEBHOOK_SECRET>"
  }
}
 ```

- To see how to use the REST API, see [API documentation](api-doc.yaml).
- Token is required to enable GitHub integration. How to get GitHub token? See [official docs](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token).
- GitHub webhook
  - To correctly handle case when git tag (used for a release in your app) is removed from the repository, you need to set up a webhook in your GitHub repository.
  - Webhook should listen to delete tag event and point it to the REST API endpoint `POST /webhooks/github/tags`.
  - If you add webhook secret to the GitHub webhook, you need to provide it in the `webhook_secret` field.
  - How to create a webhook in GitHub? See [official docs](https://docs.github.com/en/developers/webhooks-and-events/webhooks/creating-webhooks). 
  - _At current stage, the app expects that all repositories will use the same webhook secret. This will be improved in the future._

### How to enable Slack integration?

To enable Slack integration, you need to call the REST API endpoint `PATCH /organization/settings` with the following payload:

```json
{
  "slack": {
    "enabled": true,
    "token": "<SLACK_TOKEN>"
  }
}
 ```

- To see how to use the REST API, see [API documentation](api-doc.yaml).
- How to set up a Slack app and get a token? See [official docs](https://api.slack.com/authentication/token-types#bot).
- I find Slack documentation a bit confusing, so I recommend watching [this video](https://youtu.be/h94FK8h1OJU?si=J03awkzGM5VnTwMJ&t=85) to get a better understanding of how to set up a Slack app.
