# Audio Guide Bot: Production administration

Using [docker-compose](https://docs.docker.com/compose/) you can setup administration environment for you Guide Bot deployed in [Google Cloud Platform](https://cloud.google.com/).

Environment setup:
1. Save your GCP service account key, created during [deployment setup](../README.md#production-deployment), to `sa-key.json` file in this folder
0. Copy your Cloud SQL instance connection name to `INSTANCE_CONNECTION_NAME` variable in `.env` file in in this folder
0. Install [Docker](https://docs.docker.com/get-docker/) 

To start the environment use:
```sh
docker compose up -d
```

Cloud SQL database will be available from `postgres-admin` container at `cloud-sql-proxy:5432`.