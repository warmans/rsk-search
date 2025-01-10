package archive

import (
	"encoding/json"
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"os"
	"path"
	"strings"
	"sync"
)

func NewStore(archiveDir string) *Store {
	return &Store{archiveDir: archiveDir, validFiles: map[string]struct{}{}}
}

type Store struct {
	archiveDir string
	validFiles map[string]struct{}
	lock       sync.RWMutex
}

func (s *Store) IsValidFile(name string) (bool, error) {
	// has the file already been seen?
	s.lock.RLock()
	if _, ok := s.validFiles[name]; ok {
		s.lock.RUnlock()
		return ok, nil
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()

	// nope - check on disk
	archive, err := s.ListItems()
	if err != nil {
		return false, err
	}
	found := false
	for _, v := range archive {
		for _, f := range v.Files {
			s.validFiles[f] = struct{}{}
			if f == name {
				found = true
			}
		}
	}
	return found, nil
}

func (s *Store) ListItems() (models.ArchiveMetaList, error) {
	files, err := os.ReadDir(s.archiveDir)
	if err != nil {
		return nil, err
	}

	out := make(models.ArchiveMetaList, 0)
	for _, v := range files {
		if v.IsDir() || !strings.HasSuffix(v.Name(), ".meta.json") {
			continue
		}
		raw, err := os.ReadFile(path.Join(s.archiveDir, v.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read archive meta %s: %w", v.Name(), err)
		}
		meta := models.ArchiveMeta{}
		if err := json.Unmarshal(raw, &meta); err != nil {
			return nil, fmt.Errorf("failed to decode archive meta %s: %w", v.Name(), err)
		}
		out = append(out, meta)
	}

	return out, nil
}
