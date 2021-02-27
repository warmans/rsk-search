package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func WithFile(path string, cb func(f *os.File) error) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err != nil {
			if err != nil {
				err = fmt.Errorf("failed to close with error %s after error: %w", closeErr.Error(), err)
			}
			err = closeErr
		}
	}()
	err = cb(f)
	return
}

func WithJSONFile(path string, cb func(encoder *json.Encoder) error) error {
	return WithFile(path, func(f *os.File) error {
		enc := json.NewEncoder(f)
		enc.SetIndent("  ", "  ")
		return cb(enc)
	})
}
