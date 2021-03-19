package models

import "github.com/warmans/rsk-search/gen/api"

const EndSecondEOF = -1

type Tscript struct {
	Publication string  `json:"publication"`
	Series      int32   `json:"series"`
	Episode     int32   `json:"episode"`
	Chunks      []Chunk `json:"chunks"`
}

func (i Tscript) ID() string {
	return IncompleteTranscriptID(i)
}

type Chunk struct {
	ID          string `json:"id"`
	Raw         string `json:"raw"`
	StartSecond int64  `json:"start_second"`
	EndSecond   int64  `json:"end_second"`
}

func (c *Chunk) Proto(tscriptID string) *api.TscriptChunk {
	if c == nil {
		return nil
	}
	return &api.TscriptChunk{
		Id:           c.ID,
		TscriptId:    tscriptID,
		Raw:          c.Raw,
		AudioClipUri: "", //todo
	}
}

type ChunkStats struct {
	NextChunk       string
	TotalChunks     int32
	ApprovedChunks  int32
	SubmittedChunks int32
}

func (c *ChunkStats) Proto() *api.ChunkStats {
	if c == nil {
		return nil
	}
	return &api.ChunkStats{
		NumPending:           c.TotalChunks - c.ApprovedChunks,
		NumSubmitted:         c.SubmittedChunks,
		SuggestedNextChunkId: c.NextChunk,
	}
}

type Contribution struct {
	ID            string
	AuthorID      string
	ChunkID       string
	Transcription string
}
