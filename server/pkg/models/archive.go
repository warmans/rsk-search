package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"time"
)

type ArchiveMeta struct {
	OriginalMessageID string    `json:"original_message_id"`
	CreatedAt         time.Time `json:"created_at"`
	Files             []string  `json:"files"`
	Description       string    `json:"description"`
	Episode           string    `json:"episode"`
}

func (a *ArchiveMeta) Proto() *api.Archive {
	return &api.Archive{
		Id:             a.OriginalMessageID,
		Description:    a.Description,
		RelatedEpisode: a.Episode,
		Files:          a.Files,
	}
}

type ArchiveMetaList []ArchiveMeta

func (l ArchiveMetaList) Proto() *api.ArchiveList {
	out := make([]*api.Archive, len(l))
	for k, v := range l {
		out[k] = v.Proto()
	}
	return &api.ArchiveList{Items: out}
}
