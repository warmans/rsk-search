echo "Removing old files"
rm -rf sprite/*;

echo "Generating sprites"
for ep in ../data/episodes/*; do grep '"media_type": "video"' ${ep} && montage $(find -type f -name "$(basename ${ep} .json)*" | sort -V) -mode Concatenate -tile x1 sprite/$(basename ${ep} .json).png; done;

echo "creating JPGs"
for s in sprite/*.png; do echo $(basename ${s} .png); convert ${s} sprite/$(basename ${s} .png).jpg; done;
