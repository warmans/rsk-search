package spotify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// result of
type EpisodeList struct {
	Items []Episode `json:"items"`
	Next  *string   `json:"next"`
}

type Episode struct {
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

type Artist struct {
	Name string
	URI  string
}

type Track struct {
	Artists []Artist

	AlbumName string
	AlbumURI  string

	Name     string
	TrackURI string
}

func NewSearch(token string) *Search {
	return &Search{client: http.DefaultClient, token: token}
}

type Search struct {
	client *http.Client
	token  string
}

func (s *Search) FindTrack(term string) (*Track, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`https://api.spotify.com/v1/search?q=%s&type=track`, url.QueryEscape(term)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &searchResult{}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(result); err != nil {
		return nil, err
	}
	if result.Tracks == nil || len(result.Tracks.Items) == 0 {
		return nil, nil
	}

	bestMatch := result.Tracks.Items[0]

	track := &Track{
		Name:     bestMatch.Name,
		TrackURI: bestMatch.URI,
	}

	if bestMatch.Album != nil {
		track.AlbumName = bestMatch.Album.Name
		track.AlbumURI = bestMatch.Album.URI
	}

	for _, a := range bestMatch.Artists {
		track.Artists = append(track.Artists, Artist{Name: a.Name, URI: a.URI})
	}

	return track, nil
}

type searchResult struct {
	Tracks *tracks `json:"tracks"`
}
type tracks struct {
	Items []item `json:"items"`
}

type item struct {
	Album   *album   `json:"album"`
	Artists []artist `json:"artists"`
	Name    string   `json:"name"`
	URI     string   `json:"uri"`
}

type album struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type artist struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}
