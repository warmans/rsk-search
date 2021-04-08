package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"time"
)

func ContributionStateFromProto(state api.ContributionState) ContributionState {
	switch state {
	case api.ContributionState_STATE_PENDING:
		return ContributionStatePending
	case api.ContributionState_STATE_REQUEST_APPROVAL:
		return ContributionStateApprovalRequested
	case api.ContributionState_STATE_APPROVED:
		return ContributionStateApproved
	case api.ContributionState_STATE_REJECTED:
		return ContributionStateRejected
	}
	return ContributionStateUnknown
}

type ContributionState string

func (s ContributionState) Proto() api.ContributionState {
	switch s {
	case ContributionStatePending:
		return api.ContributionState_STATE_PENDING
	case ContributionStateApprovalRequested:
		return api.ContributionState_STATE_REQUEST_APPROVAL
	case ContributionStateApproved:
		return api.ContributionState_STATE_APPROVED
	case ContributionStateRejected:
		return api.ContributionState_STATE_REJECTED
	}
	return api.ContributionState_STATE_UNDEFINED
}

const (
	ContributionStatePending           ContributionState = "pending"
	ContributionStateApprovalRequested ContributionState = "request_approval"
	ContributionStateApproved          ContributionState = "approved"
	ContributionStateRejected          ContributionState = "rejected"
	ContributionStateUnknown           ContributionState = "unknown"
)

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

type TscriptStats struct {
	ID                              string
	Publication                     string
	Series                          int32
	Episode                         int32
	ChunkContributionStates         map[string][]ContributionState
	NumChunks                       int32
	NumContributions                int32
	NumPendingContributions         int32
	NumRequestApprovalContributions int32
	NumApprovedContributions        int32
	NumRejectedContributions        int32
}

func (c *TscriptStats) AsEpisode() *Episode {
	return &Episode{
		Publication: c.Publication,
		Series:      c.Series,
		Episode:     c.Episode,
	}
}

func (c *TscriptStats) Proto() *api.TscriptStats {
	if c == nil {
		return nil
	}
	res := &api.TscriptStats{
		Id:                              c.ID,
		Publication:                     c.Publication,
		Series:                          c.Series,
		Episode:                         c.Episode,
		ChunkContributions:              map[string]*api.ChunkStates{},
		NumChunks:                       c.NumChunks,
		NumContributions:                c.NumContributions,
		NumPendingContributions:         c.NumPendingContributions,
		NumRequestApprovalContributions: c.NumRequestApprovalContributions,
		NumApprovedContributions:        c.NumApprovedContributions,
		NumRejectedContributions:        c.NumRejectedContributions,
	}

	for k, v := range c.ChunkContributionStates {
		res.ChunkContributions[k] = &api.ChunkStates{
			States: make([]api.ContributionState, len(v)),
		}
		for staK, staV := range v {
			res.ChunkContributions[k].States[staK] = staV.Proto()
		}
	}
	return res
}

type Chunk struct {
	ID          string `json:"id"`
	Raw         string `json:"raw"`
	StartSecond int64  `json:"start_second"`
	EndSecond   int64  `json:"end_second"`
}

func (c *Chunk) Proto(tscriptID string, contribCount int32) *api.TscriptChunk {
	if c == nil {
		return nil
	}
	return &api.TscriptChunk{
		Id:               c.ID,
		TscriptId:        tscriptID,
		Raw:              c.Raw,
		NumContributions: contribCount,
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
	State         ContributionState
	StateComment  string
}

func (c *Contribution) Proto() *api.ChunkContribution {
	if c == nil {
		return nil
	}
	return &api.ChunkContribution{
		Id:         c.ID,
		ChunkId:    c.ChunkID,
		Transcript: c.Transcription,
		AuthorId:   c.AuthorID,
		State:      c.State.Proto(),
	}
}
func (c *Contribution) ShortProto() *api.ShortChunkContribution {
	if c == nil {
		return nil
	}
	return &api.ShortChunkContribution{
		Id:       c.ID,
		ChunkId:  c.ChunkID,
		AuthorId: c.AuthorID,
		State:    c.State.Proto(),
	}
}

type ContributionActivity struct {
	ChunkID     string
	AccessedAt  *time.Time
	SubmittedAt *time.Time
	ApprovedAt  *time.Time
	RejectedAt  *time.Time
}

type TimelineEvent struct {
	Who  string     `db:"who"`
	What string     `db:"what"`
	When *time.Time `db:"when"`
}

func (t *TimelineEvent) Proto() *api.TscriptTimelineEvent {
	if t == nil {
		return nil
	}
	res := &api.TscriptTimelineEvent{
		Who:  t.Who,
		What: t.What,
	}
	if t.When != nil {
		res.When = t.When.Format(time.RFC3339)
	}
	return res
}
