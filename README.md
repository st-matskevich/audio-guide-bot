# Audio Guide Bot
Telegram bot for taking audio tours. Built on top of [Telegram Mini App Template](https://github.com/st-matskevich/tg-mini-app-template). Available at [@audio_guide_bot](https://t.me/audio_guide_bot).

## Setup prerequisites
[Telegram Bot](https://core.telegram.org/bots) token is required to interact with [Telegram Bot API](https://core.telegram.org/bots/api). To get one, сreate a bot using [@BotFather](https://t.me/botfather) or follow [Telegram instructions](https://core.telegram.org/bots#how-do-i-create-a-bot).

[Telegram Payments](https://core.telegram.org/bots/payments) token is required to start accepting payments for tickets. To get one, request it from [@BotFather](https://t.me/botfather) or follow [Telegram instructions](https://core.telegram.org/bots/payments#connecting-payments).

## Local environment
This repository provides an easy-to-use local development environment. Using it you can start writing your bot business logic without spending time on the environment.

Local environment includes:
- [ngrok](https://ngrok.com/) reverse proxy to server local mini-app and bot deployment over HTTPS
- [nginx](https://www.nginx.com/) reverse proxy to host both API and UI on one ngrok domain and thus fit into the [free plan](https://ngrok.com/pricing)
- React fast refresh to avoid rebuilding docker container on each change of the UI code
- [PostgreSQL](https://www.postgresql.org/) database instance and [pgAdmin](https://github.com/pgadmin-org/pgadmin4) instance to manage it
- [s3gw](https://github.com/aquarist-labs/s3gw) S3 instance and [s3gw-ui](https://github.com/aquarist-labs/s3gw-ui) instance to manage it

Local environment setup:
1. Create an account on [ngrok](https://ngrok.com/)
0. Get a [ngrok auth token](https://ngrok.com/docs/secure-tunnels/ngrok-agent/tunnel-authtokens/) and save it to `NGROK_AUTHTOKEN` variable in `.env` file in the project root directory
0. Claim a [free ngrok domain](https://ngrok.com/blog-post/free-static-domains-ngrok-users) and save it to `NGROK_DOMAIN` variable in `.env` file in the project root directory
0. Copy [Telegram Bot token](#setup-prerequisites) and save it to `TELEGRAM_BOT_TOKEN` variable in `.env` file in the project root directory
0. Copy [Telegram Payments token](#setup-prerequisites) and save it to `TELEGRAM_PAYMENTS_TOKEN` variable in `.env` file in the project root directory
0. Generate random string for JWT signing secret and save it to `JWT_SECRET` variable in `.env` file in the project root directory
0. Install [Docker](https://docs.docker.com/get-docker/)

To start or update the environment with the latest code changes, use:
```sh
docker compose up --build -d
```

After successful deployment, your local bot API will be available at https://ngrok-domain/api. Use this URL to set the bot webhook as described [switching bot environment](#switching-bot-environment).

## Production deployment
This repository provides a [workflow](https://docs.github.com/actions) to automatically deploy the code to [Google Cloud Platform](https://cloud.google.com/). The deploy job is triggered on each push to the [main](https://github.com/st-matskevich/audio-guide-bot/tree/main) branch.

GCP services used for deployment:
- [Cloud Run](https://cloud.google.com/run) to host dockerized API and UI code
- [Artifact Registry](https://cloud.google.com/artifact-registry) to store docker images
- [Secret Manager](https://cloud.google.com/secret-manager) to store sensitive data
- [Cloud SQL](https://cloud.google.com/sql) to host relational database
- [Cloud Storage](https://cloud.google.com/storage) to store binary data

Deployment setup:
1. [Create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project) in GCP
0. Copy project ID and save it to `GCP_PROJECT_ID` GitHub variable
0. [Pick a region](https://cloud.withgoogle.com/region-picker/) for your app and save it to `GCP_PROJECT_REGION` GitHub variable
0. [Create a service account](https://cloud.google.com/iam/docs/service-accounts-create#creating) with the following rights:
    - Service Account User (to create resources by the name of this account)
    - Cloud Run Admin (to create Cloud Run instances)
    - Artifact Registry Administrator (to manage images in the registry)
    - Secret Manager Secret Accessor (to access GCP secrets)
    - Cloud SQL Client (to connect to the Cloud SQL instance)
    - Storage Object User (to access objects in Cloud Storage bucket)
0. Copy the service account email and save it to `GCP_SA_EMAIL` GitHub variable
0. [Export the service account key](https://cloud.google.com/iam/docs/keys-create-delete#creating) and save it to `GCP_SA_KEY` GitHub secret
0. Enable the following GCP APIs:
    - Cloud Run Admin API (to create Cloud Run instances)
    - Secret Manager API (to securely store secrets)
    - Compute Engine API (to create Cloud SQL instances)
    - Cloud SQL Admin API (to connect Cloud Run to Cloud SQL)
0. Create [Artifact Registry for Docker images](https://cloud.google.com/artifact-registry/docs/docker/store-docker-container-images#create) in `GCP_PROJECT_REGION` region
0. Copy Artifact Registry name and save it to `GCP_ARTIFACT_REGISTRY` GitHub variable
0. [Create a PostgreSQL Cloud SQL instance](https://cloud.google.com/sql/docs/postgres/create-instance#create-2nd-gen) in `GCP_PROJECT_REGION` region
    - If you want to reduce the cost of the instance:
      - You can implement [instance scheduler](https://cloud.google.com/blog/topics/developers-practitioners/lower-development-costs-schedule-cloud-sql-instances-start-and-stop)
      - You can set "Zonal availability" to "Single zone" instead of "Multiple zones (Highly available)"
      - You can set machine configuration to the minimal one (1 shared vCPU, 0.614 GB RAM)
      - You can set storage type to HDD
      - You can reduce storage size to 10 GB
0. Copy Cloud SQL instance connection name, which can be found on the Overview page for your instance, and save it to `GCP_SQL_INSTANCE_CONNECTION_NAME` GitHub variable
0. [Create a user](https://cloud.google.com/sql/docs/postgres/create-manage-users#creating) for your Cloud SQL instance
0. [Create a database](https://cloud.google.com/sql/docs/postgres/create-manage-databases#create) for your Cloud SQL instance
0. [Create a Cloud Storage bucket](https://cloud.google.com/storage/docs/creating-buckets#create_a_new_bucket) in `GCP_PROJECT_REGION` region
0. [Create an HMAC key](https://cloud.google.com/storage/docs/authentication/managing-hmackeys#create) for your service account
0. [Create the following secrets](https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets#create) in Secret Manager:
    - [Telegram Bot token](#setup-prerequisites) and save the secret name to `GCP_SECRET_TG_BOT_TOKEN` GitHub variable
    - [Telegram Payments token](#setup-prerequisites) and save the secret name to `GCP_SECRET_TG_PAYMENTS_TOKEN` GitHub variable
    - Random string for JWT signing secret and save the secret name to `GCP_SECRET_JWT_SECRET` GitHub variable
    - PostgreSQL connection string and save the secret name to `GCP_SECRET_DB_URL` GitHub variable
      - Connection string format is `postgres://{USER}:{PASSWORD}@/{DATABASE}?host={HOST}`, where:
        - `{USER}` is the name of the user created for Cloud SQL instance above
        - `{PASSWORD}` is the password of the user created for Cloud SQL instance above
        - `{HOST}` is the `/cloudsql/{INSTANCE_CONNECTION_NAME}`, where `{INSTANCE_CONNECTION_NAME}` is the Cloud SQL instance connection name, which is equal to `GCP_SQL_INSTANCE_CONNECTION_NAME`
        - `{DATABASE}` is the name of the database created for Cloud SQL instance above
    - Cloud Storage connection string and save the secret name to `GCP_SECRET_S3_URL` GitHub variable
      - Connection string format is `https://{KEY}:{SECRET}@storage.googleapis.com/{BUCKET}`, where:
        - `{KEY}` is the key of the service account HMAC key created above
        - `{SECRET}` is the secret of the service account HMAC key created above
        - `{BUCKET}` is the name of the bucket created in the Cloud Storage above
0. Define the following GitHub variables:
    - `GCP_SERVICE_MIGRATOR_NAME` with the desired name of Migrator Cloud Run instance 
    - `GCP_SERVICE_UI_NAME` with the desired name of UI Cloud Run instance 
    - `GCP_SERVICE_UI_MAX_INSTANCES` with the desired maximum number of UI service instances
    - `GCP_SERVICE_API_NAME` with the desired name of API Cloud Run instance 
    - `GCP_SERVICE_API_MAX_INSTANCES` with the desired maximum number of API service instances   

After successful deployment, obtain the bot API URL from either `deploy-api` job results or from [GCP Project Console](https://console.cloud.google.com) and proceed to [switching bot environment](#switching-bot-environment).

## Switching bot environment
After the bot is either [launched locally](#local-environment) or [deployed in GCP](#production-deployment), Telegram needs to be configured with a proper webhook URL. To set it, use:
```sh
curl https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook?url=${BOT_API_URL}/bot
```

## Project structure
- Project root directory - contains files for [local environment](#local-environment)
  - [docker-compose.yml](docker-compose.yml) - docker-compose file to setup  [local environment](#local-environment)
  - [proxy.template](proxy.template) - nginx config template to route ngrok domain to API and UI containers
- [api](/api/) - contains files for the bot backend written in Go
  - [Dockerfile](/api/Dockerfile) - docker file for the API service image
  - Other files - minimal Go package layout with HTTP endpoint to handle [Telegram Bot API](https://core.telegram.org/bots/api) calls
- [ui](/ui/) - contains files for the mini app UI written in JS with React
  - [dev.dockerfile](/ui/dev.dockerfile) - docker file for the [local environment](#local-environment) to track code changes and updated served UI bundle
  - [prod.dockerfile](/ui/prod.dockerfile) - docker file for the [production deployment](#production-deployment) to bundle UI and serve via nginx
  - [nginx.template](/ui/nginx.template) - nginx config template for production UI service image
  - Other files - React template with minimal changes for the integration with Telegram

## Built with
- [Docker](https://www.docker.com/)
- [Go](https://go.dev/)
- [React](https://react.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Fiber](https://github.com/gofiber/fiber)
- [gotgbot](https://github.com/PaulSonOfLars/gotgbot)
- [minio-go](https://github.com/minio/minio-go)
- [migrate](https://github.com/golang-migrate/migrate)
- [s3gw](https://github.com/aquarist-labs/s3gw)
- [s3gw-ui](https://github.com/aquarist-labs/s3gw-ui)
- [pgAdmin](https://github.com/pgadmin-org/pgadmin4)
- [nginx](https://www.nginx.com/)
- [ngrok](https://ngrok.com/)

## License
Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## Contributing
Want a new feature added? Found a bug?
Go ahead and open [a new issue](https://github.com/st-matskevich/audio-guide-bot/issues/new) or feel free to submit a pull request.