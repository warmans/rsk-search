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

## Importing video transcripts

Rename the files first to make life easier e.g.

```
for ep in *; do mv "${ep}" "example-$(echo ${ep} | awk '{print $5}').mp4"; done;
```

### 1. If the videos are mp4s with embedded subtitles extract the subs: 

```
for i in $(seq 1 6); do ffmpeg -i orig/example-S01E0${i}.mkv example-S01E0${i}.srt; done;
```

### 2. Resize video so gifs need less processing

```
for i in $(seq 1 6); do ffmpeg -i orig/example-S01E0${i}.mp4 -filter_complex "[0:v]fps=10,scale=598:-1" example-S01E0${i}.mp4; done;
```

Note that this can fail due to the resolution not being divisible by two. Just change it slightly and retry.

### 3. Extract create transcripts from the subs: 

Note that BOMs are stripped just in case with sed.

```
 for i in $(seq 1 6); do sed -i '1s/^\xEF\xBB\xBF//' path/to/example-S01E0${i}.srt; ./bin/rsk-search data init-from-srt --srt-path path/to/example-S01E0${i}.srt -p example -s 1 -e ${i} -m path/to/example-S01E0${i}.mp4; done
```
