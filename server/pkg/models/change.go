package models

import (
	"time"

	"github.com/warmans/rsk-search/gen/api"
)

type TranscriptChangeCreate struct {
	AuthorID          string
	EpID              string
	TranscriptVersion string
	Name              string
	Summary           string
	ReleaseDate       time.Time
	Transcription     string
}

type TranscriptChangeUpdate struct {
	ID            string
	Name          string
	Summary       string
	ReleaseDate   time.Time
	Transcription string
	State         ContributionState
}

type TranscriptChange struct {
	ID                string
	EpID              string
	TranscriptVersion string
	Author            *Author
	Name              string
	Summary           string
	ReleaseDate       time.Time
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
		Name:              c.Name,
		Summary:           c.Summary,
		ReleaseDate:       c.ReleaseDate.Format(time.DateOnly),
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
