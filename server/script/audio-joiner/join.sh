echo "Creating manifest..."
for dir in `ls | grep 'series-*'`; do for f in `ls ${dir}/* | grep .mp3`; do echo "file '${f}'"; done; done > manifest.txt

echo "Creating xfm.mp3..."
ffmpeg -f concat -safe 0 -i manifest.txt -c copy xfm.mp3
