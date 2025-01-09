package models

type ArchiveMeta struct {
	OriginalMessageID string   `json:"original_message_id"`
	Files             []string `json:"files"`
	Description       string   `json:"description"`
	Episode           string   `json:"episode"`
}
