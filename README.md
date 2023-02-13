# rsk-search

## Setup 

## GUI Development

From the gui directory:

1. Install dependencies `npm install`.
2. Start the development server with `npm run start-prod`. This will not require any 
local server and will use the live API instead.

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

### How to change the API 

1. Edit proto file e.g. `proto/search.proto`.
2. Run `make generate`.
3. Update code to reflect changes in proto file e.g. `pkg/service/grpc/search.go`.
4. In `gui` directory run `npm run generate-api-client` to sync the GUI client with the latest API definitions. 

### How to update transcripts

1. Update data in `./var/data/episodes`
2. Run `make update.transcriptions SPOTIFY_TOKEN=[my token]`
3. Re-generate DB/Search index with `make init.all`.

Generating a spotify token. You can just generate one here if you are logged in:
https://developer.spotify.com/console/get-search-item/ (click GET TOKEN)

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
