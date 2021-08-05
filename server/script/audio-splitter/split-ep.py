import argparse
import json

from pydub import AudioSegment

parser = argparse.ArgumentParser(description='Split audio file into chunks according to chunked transcription')
parser.add_argument('--meta', type=str,
                    help='the chunked transcript')
parser.add_argument('--audio', type=str,
                    help='the corresponding audio file')
parser.add_argument('--outpath', type=str,
                    help='dir to save output files')

args = parser.parse_args()

with open(args.meta) as f:
    data = json.load(f)

episode = AudioSegment.from_wav(args.audio)

for chunk in data["chunks"]:
    print("chunk", chunk["id"], chunk["start_second"], chunk["end_second"])
    if chunk["end_second"] == -1:
        audioChunk = episode[chunk["start_second"] * 1000:]
    else:
        audioChunk = episode[chunk["start_second"] * 1000:chunk["end_second"] * 1000]

    audioChunk.export(args.outpath + "/" + chunk["id"] + ".mp3", format="mp3")
