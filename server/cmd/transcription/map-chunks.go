package transcription

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/transcript"
	"github.com/warmans/rsk-search/pkg/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const targetChunkDuration = 180

func MapChunksCmd() *cobra.Command {

	var incompleteFilesDir string

	cmd := &cobra.Command{
		Use:   "map-chunks",
		Short: "Take raw transcript and map it to a series of chunks that can be consumed by the API",
		RunE: func(cmd *cobra.Command, args []string) error {

			dirEntries, err := ioutil.ReadDir(incompleteFilesDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {
				if dirEntry.IsDir() {
					continue
				}
				if strings.HasPrefix(dirEntry.Name(), ".") {
					continue
				}
				f, err := os.Open(path.Join(incompleteFilesDir, dirEntry.Name()))
				if err != nil {
					return err
				}
				ts, err := getIncompleteTranscriptionModelFromName(dirEntry.Name())
				if err != nil {
					return err
				}
				ts.Chunks, err = getChunks(bufio.NewScanner(f))
				if err != nil {
					return err
				}

				outFile := path.Join(incompleteFilesDir, "..", "chunked", dirEntry.Name())

				exists, err := util.FileExists(outFile)
				if err != nil {
					return err
				}
				if exists {
					fmt.Printf("skipping already mapped file: %s\n", outFile)
					continue
				}
				if err := util.WithCreateJSONFileEncoder(outFile, func(enc *json.Encoder) error {
					return enc.Encode(ts)

				}); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&incompleteFilesDir, "incomplete-dir", "d", "./var/data/incomplete/raw", "Path to incomplete transcripts created by gcloud")

	return cmd
}

func getChunks(scanner *bufio.Scanner) ([]models.Chunk, error) {
	chs := make([]models.Chunk, 0)
	var curentChunk *models.Chunk

	for scanner.Scan() {
		line := scanner.Text()

		if transcript.IsOffsetTag(line) {
			offsetSeconds, ok := transcript.ScanOffset(line)
			if !ok {
				return nil, fmt.Errorf("failed to get valid offset from line: %s", line)
			}
			if curentChunk == nil {
				curentChunk = &models.Chunk{ID: shortuuid.New(), StartSecond: offsetSeconds}
			}
			if offsetSeconds - curentChunk.StartSecond >= targetChunkDuration {
				curentChunk.EndSecond = offsetSeconds
				chs = append(chs, *curentChunk)
				curentChunk = &models.Chunk{ID: shortuuid.New(), StartSecond: offsetSeconds}
			}
		} else {
			if curentChunk == nil {
				return nil, fmt.Errorf("file seems to be missing initial offset")
			}
		}
		curentChunk.Raw += line + "\n"
	}

	if curentChunk != nil && len(curentChunk.Raw) > 0 {
		curentChunk.EndSecond = models.EndSecondEOF
		chs = append(chs, *curentChunk)
	}
	return chs, scanner.Err()
}

func getIncompleteTranscriptionModelFromName(fileName string) (*models.Tscript, error) {

	fileName = strings.TrimSuffix(fileName, ".txt")

	publicationAndSeries := strings.Split(fileName, "-")
	if len(publicationAndSeries) != 2 {
		return nil, fmt.Errorf("could not parse publication from filename: %s", fileName)
	}
	series, episode, err := models.ParseStandardEpisodeName(publicationAndSeries[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse series/episode from filename: %s", publicationAndSeries[1])
	}
	return &models.Tscript{
		Publication: publicationAndSeries[0],
		Series:      series,
		Episode:     episode,
	}, nil
}
