#!/usr/bin/env bash

for file in *-S0*E*.mp3; do
  if [[ $file =~ ^(.+)-S0([0-9]+)(E[0-9]+)\.mp3$ ]]; then
    prefix="${BASH_REMATCH[1]}"
    season="${BASH_REMATCH[2]}"
    rest="${BASH_REMATCH[3]}"

    newfile="${prefix}-S${season}${rest}.mp3"

    if [[ "$file" != "$newfile" ]]; then
      mv -- "$file" "$newfile"
      echo "Renamed: $file -> $newfile"
    fi
  fi
done