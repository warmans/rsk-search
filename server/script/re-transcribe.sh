EP=S1E23; ./bin/rsk-search data transcribe-assembly-ai -i "https://scrimpton.com/dl/media/xfm-${EP}.mp3?remastered=1" && \
 ./bin/rsk-search data aai-merge-timestamps -s "var/aai-transcripts/xfm-${EP}.mp3?remastered=1.json" -t ep-xfm-${EP}.json && \
 ./bin/rsk-search data aai-merge-sentiments -s "var/aai-transcripts/xfm-${EP}.mp3?remastered=1.json" -t "ep-xfm-${EP}.json" && \
 ./bin/rsk-search data promote-remaster-audio -s "xfm-${EP}"
