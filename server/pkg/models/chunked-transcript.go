package models

import (
	"encoding/json"
	"fmt"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/util"
	"path"
	"strings"
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

type ChunkedTranscript struct {
	Publication string  `json:"publication"`
	Series      int32   `json:"series"`
	Episode     int32   `json:"episode"`
	Name        string  `json:"name"`
	Chunks      []Chunk `json:"chunks"`
}

func (t ChunkedTranscript) ID() string {
	return IncompleteTranscriptID(t)
}

type ChunkedTranscriptStats struct {
	ID                              string
	Publication                     string
	Series                          int32
	Episode                         int32
	Name                            string
	ChunkContributionStates         map[string][]ContributionState
	NumChunks                       int32
	NumContributions                int32
	NumPendingContributions         int32
	NumRequestApprovalContributions int32
	NumApprovedContributions        int32
	NumRejectedContributions        int32
}

func (c *ChunkedTranscriptStats) AsEpisode() *Transcript {
	return &Transcript{
		Publication: c.Publication,
		Series:      c.Series,
		Episode:     c.Episode,
		Name:        c.Name,
	}
}

func (c *ChunkedTranscriptStats) Proto() *api.ChunkedTranscriptStats {
	if c == nil {
		return nil
	}
	res := &api.ChunkedTranscriptStats{
		Id:                              c.ID,
		Publication:                     c.Publication,
		Series:                          c.Series,
		Episode:                         c.Episode,
		Name:                            c.Name,
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
	ID               string        `json:"id" db:"id"`
	TscriptID        string        `json:"tscript_id" db:"tscript_id"`
	Raw              string        `json:"raw" db:"raw"`
	StartSecond      time.Duration `json:"start_second" db:"start_second"`
	EndSecond        time.Duration `json:"end_second" db:"end_second"`
	NumContributions int32         `json:"num_contributions" db:"num_contributions"`
}

func (c *Chunk) Proto() *api.Chunk {
	if c == nil {
		return nil
	}
	return &api.Chunk{
		Id:                  c.ID,
		ChunkedTranscriptId: c.TscriptID,
		Raw:                 c.Raw,
		NumContributions:    c.NumContributions,
		// this is hacky. It's annoying, the DB should really just have the episode id.
		// the good thing is the transcript IS the episode id with a different prefix...
		EpisodeId:   strings.Replace(c.TscriptID, "ts-", "ep-", 1),
		StartTimeMs: int32(c.StartSecond.Milliseconds()),
		EndTimeMs:   int32(c.EndSecond.Milliseconds()),
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

type ContributionCreate struct {
	AuthorID      string
	ChunkID       string
	Transcription string
	State         ContributionState
}

type ContributionUpdate struct {
	ID            string
	AuthorID      string
	Transcription string
	State         ContributionState
}

type ChunkContribution struct {
	ID            string
	TscriptID     string
	ChunkID       string
	Author        *ShortAuthor
	Transcription string
	State         ContributionState
	StateComment  string
	CreatedAt     time.Time
}

func (c *ChunkContribution) Proto() *api.ChunkContribution {
	if c == nil {
		return nil
	}
	return &api.ChunkContribution{
		Id:           c.ID,
		ChunkId:      c.ChunkID,
		Transcript:   c.Transcription,
		State:        c.State.Proto(),
		StateComment: c.StateComment,
		Author:       c.Author.Proto(),
		CreatedAt:    c.CreatedAt.Format(time.RFC3339),
	}
}

func (c *ChunkContribution) ShortProto() *api.ShortChunkContribution {
	if c == nil {
		return nil
	}
	return &api.ShortChunkContribution{
		Id:       c.ID,
		ChunkId:  c.ChunkID,
		AuthorId: c.Author.ID,
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

type TscriptImportCreate struct {
	EpID   string `json:"epid" db:"epid"`
	EpName string `json:"epname" db:"epname"`
	Mp3URI string `json:"mp3_uri" db:"mp3_uri"`
}

type TscriptImport struct {
	ID          string     `json:"id" db:"id"`
	EpID        string     `json:"epid" db:"epid"`
	EpName      string     `json:"epname" db:"epname"`
	Mp3URI      string     `json:"mp3_uri" db:"mp3_uri"`
	Log         string     `json:"-" db:"log"`
	CreatedAt   *time.Time `json:"-" db:"created_at"`
	CompletedAt *time.Time `json:"-" db:"completed_at"`
}

func (c *TscriptImport) Proto() *api.TscriptImport {
	if c == nil {
		return nil
	}
	logs := make([]*api.TscriptImportLog, 0)
	if c.Log != "" {
		if err := json.Unmarshal([]byte(c.Log), &logs); err != nil {
			// just ignore it for now. The logs aren't used for anything yet.
			logs = make([]*api.TscriptImportLog, 0)
		}
	}
	return &api.TscriptImport{
		Id:          c.ID,
		Epid:        c.EpID,
		Mp3Uri:      c.Mp3URI,
		Log:         logs,
		CreatedAt:   util.FormatTimeForRPCResponse(c.CreatedAt),
		CompletedAt: util.FormatTimeForRPCResponse(c.CompletedAt),
	}
}

func (c *TscriptImport) WorkingDir(parentDirParth string) string {
	return path.Join(parentDirParth, c.ID)
}

func (c *TscriptImport) Mp3() string {
	return fmt.Sprintf("%s.mp3", c.EpID)
}

func (c *TscriptImport) WAV() string {
	return fmt.Sprintf("%s.wav", c.EpID)
}

func (c *TscriptImport) MachineTranscript() string {
	return fmt.Sprintf("%s.machine.txt", c.EpID)
}

func (c *TscriptImport) ChunkedMachineTranscript() string {
	return fmt.Sprintf("%s.chunks.json", c.EpID)
}

type TscriptImportLog struct {
	Stage string `json:"stage"`
	Msg   string `json:"msg"`
}
