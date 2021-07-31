package main

import (
	"flag"
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
)

// rename original mp3s to a easier to work with format.
// can be converted to mp3 with:
// for f in *; do ffmpeg -i "$f" "$(basename $f .mp3).wav"; done;
func main() {

	readDir := flag.String("read-dir", "./raw", "Directory containing raw MP3s")
	publication := flag.String("publication", "xfm", "Prefix renamed files with this text")
	//wavDir := flag.String("wav-dir", "./wav-out", "Directory to write wav files")
	flag.Parse()

	files, err := os.ReadDir(*readDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		newFileName := convertRawFileName(f.Name())
		if newFileName == "" {
			log.Println("failed to convert file name: ", f.Name())
			continue
		}

		log.Println("Moving ", f.Name(), " to ", newFileName)
		if err := os.Rename(path.Join(*readDir, f.Name()), path.Join(*readDir, fmt.Sprintf("%s-%s.mp3", *publication, convertRawFileName(f.Name())))); err != nil {
			log.Fatal(err.Error())
		}
	}
}

var matchFileName = regexp.MustCompile("Series ([0-9]+) Episode ([0-9]+) .*")

// From Series 1 Episode 20 (13. April 2002).mp3 -> S1E20
func convertRawFileName(name string) string {
	matches := matchFileName.FindAllStringSubmatch(name, -1)
	if len(matches) != 1 || len(matches[0]) != 3 {
		return ""
	}
	series, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return ""
	}
	episode, err := strconv.Atoi(matches[0][2])
	if err != nil {
		return ""
	}
	return models.FormatStandardEpisodeName(int32(series), int32(episode))
}
