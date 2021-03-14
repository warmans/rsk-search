# rsk-search

## GUI Development

From the gui directory:

1. Install dependencies `npm install`.
2. Start the development server with `npm run start-prod`. This will not require any 
local server and will use the live API instead.

If the local server is running you can use `npm run start` to proxy the 
API to the local running API.

## Server development

From the server directory: 

1. Install build tools with `make install.tools`
2. Build CLI with `make vendor build`.
3. Create the DB and index with `make init.all` (only needs to be done once, or after the raw data is changed).
4. Run local server with `make run`.

### How to change the API 

1. Edit proto file in `proto/search.proto`.
2. Run `make generate`.
3. Update to reflect changes in proto file `pkg/service/grpc/search.go`.

### How to update transcripts

1. Update data in `./var/data/episodes`
2. Run `make update.transcriptions SPOTIFY_TOKEN=[my token]` (note, this can take several hours as it must re-tag all dialog)
3. Re-generate DB/Search index with `make init.all`.


Generating a spotify token. You can just generate one for the from here if you are logged in:
https://developer.spotify.com/console/get-search-item/ (click GET TOKEN)

### How to update tags

Tags can only be removed or aliased. The tag list in the meta/data dir 
is just a list of tags that the NER package will detect in the data, but 
before storing any tag against the dialog it is checked against the 
file in the tags.json file. 

1. Edit `./pkg/meta/data/tags.json`.
2. Re-tag dialog `./bin/rsk-search data transcribe` (can take several hours).
3. Re-generate DB/Search index with `make init.all`.


## Deployment

### Update scrimpton.com

(requires access to `warmans` docker hub account)

1. Create a new git tag `git tag x.x.x`.
2. `cd gui && npm run release`.
3. `cd server && make release`.
3. Re-up docker compose file on server.

### New deployment
1. Build and push docker images with your own namespace.
2. Update `./deploy/docker-compose.yaml` with correct image names/versions.
3. Copy docker-compose file to your server.
4. Run `docker-compose up -d`
