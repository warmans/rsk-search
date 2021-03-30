package data

import (
	"encoding/json"
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"os"
	"path"
)

func ReplaceEpisodeFile(dataDir string, ep *models.Episode) error {
	return util.WithReplaceJSONFileEncoder(path.Join(dataDir, fmt.Sprintf("%s.json", models.EpisodeID(ep))), func(encoder *json.Encoder) error {
		return encoder.Encode(ep)
	})
}

func SaveEpisodeToFile(dataDir string, ep *models.Episode) error {
	return util.WithCreateJSONFileEncoder(path.Join(dataDir, fmt.Sprintf("%s.json", models.EpisodeID(ep))), func(encoder *json.Encoder) error {
		return encoder.Encode(ep)
	})
}

func LoadEpisodeFile(dataDir string, publication string, name string) (*os.File, error) {
	return os.Open(path.Join(dataDir, fmt.Sprintf("ep-%s-%s.json", publication, name)))
}

func LoadEpisode(dataDir string, publication string, name string) (*models.Episode, error) {

	f, err := LoadEpisodeFile(dataDir, publication, name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	e := &models.Episode{}

	dec := json.NewDecoder(f)
	return e, dec.Decode(e)
}

func LoadEpisodePath(path string) (*models.Episode, error) {
	episode := &models.Episode{}
	if err := util.WithReadJSONFileDecoder(path, func(dec *json.Decoder) error {
		return dec.Decode(episode)
	}); err != nil {
		return nil, err
	}
	return episode, nil
}
