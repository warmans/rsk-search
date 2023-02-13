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

## GUI Development

From the gui directory:

1. Install dependencies `npm install`.
2. Start the development server with `npm run start-prod`. This will not require any 
local server and will use the live scrimpton.com API.

If the local server is running you can use `npm run start` to proxy the 
API to the local running API.

## Server development

From the server directory: 

1. Install tools `make setup`
2. Build CLI with `make vendor build`.
3. Create the DB and index with `make init.all` (only needs to be done once, or after the raw data is changed)
4. Start a local postgres instance `make dev.services.start`
5. Setup some test data `make dev.populate.chunks` 
6. Run local server with `make run`.

More info: [server README](server/README.md)

### How to change the API 

1. Edit proto file e.g. `proto/search.proto`.
2. Run `make generate`.
3. Update code to reflect changes in proto file e.g. `pkg/service/grpc/search.go`.
4. In `gui` directory run `npm run generate-api-client` to sync the GUI client with the latest API definitions. 

More info: [GUI README](gui/README.md)

## Deployment

### Update scrimpton.com

(requires access to `warmans` docker hub account)

1. Create a new git tag `git tag x.x.x`.
2. `cd gui && npm run release`.
3. `cd server && make release`.
3. Re-up docker compose file on server.

(this is done using github actions)

### New deployment

1. Build and push docker images with your own namespace.
2. Update `./deploy/docker-compose.yaml` with correct image names/versions and credentials.
3. Copy docker-compose file to your server.
4. Run `docker-compose up -d`
