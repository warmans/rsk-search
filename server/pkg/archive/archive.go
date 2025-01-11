package archive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"io"
	"net/http"
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

func (s *Store) FileExists(name string) (bool, error) {
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

func (s *Store) ArchiveFile(filename string, url string) error {
	file, err := os.OpenFile(path.Join(s.archiveDir, path.Clean(filename)), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("file already exists: %s", filename)
		}
		return fmt.Errorf("unable to archive file: internal error")
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("unable to archive file: internal error")
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to archive file: internal error")
	}

	return nil
}

func (s *Store) CreateMetadata(metadata models.ArchiveMeta) error {
	metaFile, err := s.getMetaFile(metadata.OriginalMessageID)
	if err != nil {
		return err

	}
	defer metaFile.Close()

	enc := json.NewEncoder(metaFile)
	enc.SetIndent("", "  ")
	return enc.Encode(metadata)
}

func (s *Store) getMetaFile(messageID string) (*os.File, error) {
	metaFile, err := os.OpenFile(path.Join(s.archiveDir, fmt.Sprintf("%s.meta.json", path.Clean(messageID))), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, err
		}
		// we've already stored the file, probably not worth deleting it.
		return nil, fmt.Errorf("failed to create metadata: %w", err)
	}
	return metaFile, nil
}
