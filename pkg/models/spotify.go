package models

import "time"

type SpotifyResult struct {
	Items []SpotifyItem `json:"items"`
	Next  *string       `json:"next"`
}

type SpotifyItem struct {
	AudioPreviewURL string            `json:"audio_preview_url"`
	DurationMs      int64             `json:"duration_ms"`
	ExternalUrls    map[string]string `json:"external_urls"`
	ID              string
	Name            string `json:"name"`
	ReleaseDate     string `json:"release_date"`
	URI             string `json:"uri"`

	// infer some data
	ReleaseDateT time.Time `json:"release_date_t"`
	Series       int32     `json:"series"`
	Episode      int32     `json:"episode"`
}
