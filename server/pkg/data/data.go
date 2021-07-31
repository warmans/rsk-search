package data

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"os"
	"path"
	"strings"
	"sync"
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

func LoadAllEpisodes(dataDir string) ([]*models.Episode, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}
	episodes := []*models.Episode{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		episodePath := path.Join(dataDir, entry.Name())
		ep, err := LoadEpisodePath(episodePath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load path %s", episodePath)
		}
		episodes = append(episodes, ep)
	}
	return episodes, nil
}

func NewEpisodeStore(dataDir string) (*EpisodeCache, error) {

	store := &EpisodeCache{
		cache: map[string]models.Episode{},
		lock:  sync.RWMutex{},
	}
	episodes, err := LoadAllEpisodes(dataDir)
	if err != nil {
		return nil, err
	}

	store.lock.Lock()
	defer store.lock.Unlock()
	for _, ep := range episodes {
		store.cache[ep.ID()] = *ep
	}

	return store, nil
}

var ErrNotFound = errors.New("not found")

type EpisodeCache struct {
	cache map[string]models.Episode
	lock  sync.RWMutex
}

func (s *EpisodeCache) GetEpisode(id string) (*models.Episode, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	ep, ok := s.cache[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &ep, nil
}
