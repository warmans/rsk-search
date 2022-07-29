package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"time"
)

type TranscriptChangeCreate struct {
	AuthorID          string
	EpID              string
	TranscriptVersion string
	Summary           string
	Transcription     string
}

type TranscriptChangeUpdate struct {
	ID            string
	Summary       string
	Transcription string
	State         ContributionState
}

type TranscriptChange struct {
	ID                string
	EpID              string
	TranscriptVersion string
	Author            *Author
	Summary           string
	Transcription     string
	State             ContributionState
	CreatedAt         time.Time
	Merged            bool
	PointsAwarded     float32
}

func (c *TranscriptChange) Proto() *api.TranscriptChange {
	if c == nil {
		return nil
	}
	return &api.TranscriptChange{
		Id:                c.ID,
		EpisodeId:         c.EpID,
		TranscriptVersion: c.TranscriptVersion,
		Summary:           c.Summary,
		Transcript:        c.Transcription,
		State:             c.State.Proto(),
		Author:            c.Author.ShortAuthor().Proto(),
		CreatedAt:         c.CreatedAt.Format(time.RFC3339),
		Merged:            c.Merged,
		PointsAwarded:     c.PointsAwarded,
	}
}

func (c *TranscriptChange) ShortProto() *api.ShortTranscriptChange {
	if c == nil {
		return nil
	}
	return &api.ShortTranscriptChange{
		Id:                c.ID,
		EpisodeId:         c.EpID,
		TranscriptVersion: c.TranscriptVersion,
		State:             c.State.Proto(),
		Author:            c.Author.ShortAuthor().Proto(),
		CreatedAt:         c.CreatedAt.Format(time.RFC3339),
		Merged:            c.Merged,
		PointsAwarded:     c.PointsAwarded,
	}
}
