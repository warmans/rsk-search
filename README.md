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
3. Create the DB and index with `make init.all` (only needs to be done once, or after the raw data is changed)
4. Start a local postgres instance `make dev.services.start`
5. Setup some test data `make dev.populate.chunks` 
6. Run local server with `make run`.

### How to change the API 

1. Edit proto file in `proto/search.proto`.
2. Run `make generate`.
3. Update to reflect changes in proto file `pkg/service/grpc/search.go`.
4. In `gui` directory run `npm run generate-api-client` 

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


## Transcription

To add a new episode to the pending transcriptions: 

1. Create 16bit mono .wav of the audio (if an mp3). E.g. using Audacity.
2. Upload the raw wav file to google drive.
3. Use `./bin/rsk-search transcription gcloud` command to auto-transcribe it and redirect the output into a file with a standardized name.
   * e.g. `GOOGLE_APPLICATION_CREDENTIALS=~/keys/key.json ./bin/rsk-search transcription gcloud "gs://my-bucket-name/raw/Series 4 Episode 2 (4. June 2005).wav" > ./var/data/episodes/incomplete/raw/xfm-S4E02.txt`
4. Use the `./bin/rsk-search transcription map-chunks` command to create the chunked file.
5. Create the audio chunks with the python script.
   * `python3 script/audio-splitter/split-ep.py --meta var/data/episodes/incomplete/chunked/xfm-S1E01.txt /audio/wav/xfm-S3E09.wav --outpath /out/dir`
6. Upload the chunks to google cloud.
7. Use the `./bin/rsk-search db load-tscript` command to upload the files to the DB, putting the chunks live.
