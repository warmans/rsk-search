# rsk-search

### Introduction

This is the project that powers scrimpton.com. It is split into two areas:

Backend (`/server`):

- Written in Go
- Uses Protobuff/Grpc for the API.
- Static data is stored in Sqlite, Bluge (formerly Bleve) or just flat files (see `var/data`).
- Dynamic data (e.g. in-progress contributions) is stored in Postgres.
- Audio files are stored on the server file system (no CDN or cloud storage).

Frontend (`/gui`):

- Written in Typescript (Angular)
- Styles are based on Bootstrap.
- Icons are Bootstrap Icons.

There are also some files to give an example of how the service can be deployed in `/deploy`.

## Development Prerequisites

1. tilt https://docs.tilt.dev/ (recommended)
2. docker.io, docker-compose (`apt install docker.io docker-compose`)
3. direnv (`apt install direnv` - remember to install shell hook e.g. `direnv hook bash > ~/.bashrc`)
4. npm/nodejs 18.x
5. go 1.20

## Running services manually

### GUI Development

From the gui directory:

1. Install dependencies `npm install`.
2. Start the development server with `npm run start-prod`. This will not require any
   local server and will use the live scrimpton.com API.

If the local server is running you can use `npm run start` to proxy the
API to the local running API.

### Server development

From the server directory:

Setup:

1. Install tools `make setup`
2. Build CLI with `make build`.
3. Create the DB and index with `make init.all` (only needs to be done once, or after the raw data is changed)
4. Start a local postgres instance `make dev.services.start`
6. Run local server with `make run`.

Note that invalid API keys will be used (see .envrc). These may cause some things not to work correctly.

More info: [server README](server/README.md)

#### How to change the API

1. Edit proto file e.g. `proto/search.proto`.
2. Run `make generate`.
3. Update code to reflect changes in proto file e.g. `pkg/service/grpc/search.go`.
4. In `gui` directory run `npm run generate-api-client` to sync the GUI client with the latest API definitions.

More info: [GUI README](gui/README.md)

## Running services with Tilt

Once everything is set up, you can then switch to using tilt to run all the services.

With nothing running you can run `tilt up` from the root directory, and it will start everything for the API and UI.

Tilt has its own UI to show the services' status, but you can access the scrimpton UI on the normal
port: http://localhost:4200

## Deployment

### New deployment

1. Build and push docker images with your own namespace.
2. Update `./deploy/docker-compose.yaml` with correct image names/versions and credentials.
3. Copy docker-compose file to your server.
4. Run `docker-compose up -d`
