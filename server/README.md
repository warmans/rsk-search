# Server

The important thing to know about this server is most of the data (all the important data, anyway) is stored in this repo 
as static files. The only data that lives on the server is audio files/media and temporary data like pending transcripts.

Search works by generating an index from the flat files in the deployment pipeline. When you load a transcript, it 
pretty much just returns the flat file (they're included in the docker image).

### `cmd`

Various application entrypoints including the main `server` command. If you want to understand how the API works this is a good place to 
start.

Most of the commands in here are called by the `Makefile`, which will generally give more context or ensure they're called 
in the right order.

### `gen` 

Generated files. These are updated by `make generate` and should never be manually changed.

### `pkg` 

Standalone libraries. This is the majority of the code the site uses, but they won't make much sense in isolation.

### `proto`

The proto files that define the API endpoints. These are used by `make generate`.

### `script`

Misc scripts. They are in general single-use, which is why they were not just cli commands (with some exceptions).

### `service`

This is where all the "business logic" of the API is (including the proto service implementations).

### `var`

Various files. This is where all the raw JSON for all the transcripts live.
