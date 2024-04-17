package speech2text

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/transcript"
	"io"
	"time"
)

const targetChunkDuration time.Duration = 180 * time.Second

func MapChunksFromRawTranscript(epid string, epName string, inFile io.Reader, outputWriter io.Writer) error {

	ts, err := getIncompleteTranscriptionModelFromName(epid, epName)
	if err != nil {
		return err
	}
	ts.Chunks, err = getChunks(bufio.NewScanner(inFile))
	if err != nil {
		return err
	}

	return json.NewEncoder(outputWriter).Encode(ts)
}

func getChunks(scanner *bufio.Scanner) ([]models.Chunk, error) {
	chs := make([]models.Chunk, 0)
	var curentChunk *models.Chunk

	for scanner.Scan() {
		line := scanner.Text()

		if transcript.IsOffsetTag(line) {
			offset, ok := transcript.ScanOffset(line)
			if !ok {
				return nil, fmt.Errorf("failed to get valid offset from line: %s", line)
			}
			if curentChunk == nil {
				curentChunk = &models.Chunk{ID: shortuuid.New(), StartSecond: offset}
			}
			if offset-curentChunk.StartSecond >= targetChunkDuration {
				curentChunk.EndSecond = offset
				chs = append(chs, *curentChunk)
				curentChunk = &models.Chunk{ID: shortuuid.New(), StartSecond: offset}
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

func getIncompleteTranscriptionModelFromName(epid string, epName string) (*models.ChunkedTranscript, error) {
	publication, series, episode, err := models.ParseEpID(epid)
	if err != nil {
		return nil, err
	}
	return &models.ChunkedTranscript{
		Publication: publication,
		Series:      series,
		Episode:     episode,
		Name:        epName,
	}, nil
}
