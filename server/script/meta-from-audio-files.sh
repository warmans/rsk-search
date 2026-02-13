#!/usr/bin/env bash

for f in "$@";
 do echo "Processing... ${f}" && ./bin/rsk-search data init-from-audio --audio-file-path="${f}" --publication-type=radio --publication=podcast;
done;