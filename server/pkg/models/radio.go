package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"time"
)

type RadioState struct {
	AuthorID         string        `db:"author_id"`
	EpisodeID        string        `db:"episode_id"`
	StartedAt        time.Time     `db:"started_at"`
	CurrentTimestamp time.Duration `db:"current_timestamp_ms"`
}

func (s *RadioState) Proto() *api.RadioState {
	return &api.RadioState{
		CurrentEpisode: &api.CurrentEpisode{
			ShortId:   s.EpisodeID,
			StartedAt: s.StartedAt.Format(time.RFC3339),
		},
		CurrentTimestampMs: int32(s.CurrentTimestamp.Milliseconds()),
	}
}
