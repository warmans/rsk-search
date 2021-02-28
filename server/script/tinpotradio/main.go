package main

import (
	"encoding/json"
	"fmt"
	"github.com/warmans/rsk-search/pkg/spotify"
	"github.com/warmans/rsk-search/pkg/util"
	"net/http"
	"os"
	"sort"
	"time"
)

func main() {

	spotifyToken := os.Getenv("SPOTIFY_TOKEN")
	if spotifyToken == "" {
		panic("no SPOTIFY_TOKEN in env")
	}

	episodeDateNameMap := map[string]string{}
	allItems := []spotify.Episode{}

	result := &spotify.EpisodeList{Next: stringP("https://api.spotify.com/v1/shows/2vconSkxZmWl3H2El3ZH2Q/episodes?offset=0&limit=50")}
	for result.Next != nil {
		req, err := http.NewRequest(http.MethodGet, *result.Next, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyToken))
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-type", "application/json")

		fmt.Println("fetching ", *result.Next)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		result = &spotify.EpisodeList{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(result); err != nil {
			panic(err)
		}

		for _, itm := range result.Items {

			// add a go timestamp to make it easier to parse/sort
			ts, err := time.Parse(util.ShortDateFormat, itm.ReleaseDate)
			if err != nil {
				panic(err)
			}

			// for some reason the tinpot radio release dates are offset by -1 day. I suspect something is doing
			// a timezone conversion incorrectly, but as it seems consistent it's easy to fix.
			ts = ts.AddDate(0, 0, 1)

			itm.ReleaseDateT = ts
			itm.ReleaseDate = ts.Format(util.ShortDateFormat)

			// parse name numbers
			itm.Series, itm.Episode = mustParseName(itm.Name)

			// ensure consistent naming
			itm.Name = util.FormatStandardEpisodeName(itm.Series, itm.Episode)

			allItems = append(allItems, itm)
			episodeDateNameMap[itm.ReleaseDateT.Format(time.RFC3339)] = util.FormatStandardEpisodeName(itm.Series, itm.Episode)
		}
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}

	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].ReleaseDateT.Before(allItems[j].ReleaseDateT)
	})

	// main episode meta
	if err := util.WithNewFile("./raw/xfm-spotify-meta.json", func(f *os.File) error {
		enc := json.NewEncoder(f)
		enc.SetIndent("  ", "  ")
		if err := enc.Encode(allItems); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	// also dump a map of dates -> episode names
	if err := util.WithNewFile("./raw/xfm-episode-map.json", func(f *os.File) error {
		enc := json.NewEncoder(f)
		enc.SetIndent("  ", "  ")
		if err := enc.Encode(episodeDateNameMap); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

}

func stringP(str string) *string {
	return &str
}

// S1E02
func mustParseName(name string) (int32, int32) {
	series, episode, err := util.ParseStandardEpisodeName(name)
	if err != nil {
		panic(err)
	}
	return series, episode
}
