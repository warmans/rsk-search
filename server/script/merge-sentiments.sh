#!/usr/bin/env bash

for S in {1..2}; do
  for E in {1..51};
    do EP=$(printf "S%sE%02d" "${S}" "${E}"); stat "var/aai-transcripts/xfm-${EP}.mp3?remastered=1.json" 1>/dev/null 2>/dev/null && ./bin/rsk-search data aai-merge-sentiments -s "var/aai-transcripts/xfm-${EP}.mp3?remastered=1.json" -t "ep-xfm-${EP}.json";
    done;
done;