package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func WithNewFile(path string, cb func(f *os.File) error) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close with error %s after error: %w", closeErr.Error(), err)
			}
			err = closeErr
		}
	}()
	err = cb(f)
	return
}

func WithExistingFile(path string, cb func(f *os.File) error) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("failed to close with error %s after error: %w", closeErr.Error(), err)
			}
			err = closeErr
		}
	}()
	err = cb(f)
	return
}

func WithJSONFileEncoder(path string, cb func(enc *json.Encoder) error) error {
	return WithNewFile(path, func(f *os.File) error {
		enc := json.NewEncoder(f)
		enc.SetIndent("  ", "  ")
		return cb(enc)
	})
}

func WithJSONFileDecoder(path string, cb func(dec *json.Decoder) error) error {
	return WithExistingFile(path, func(f *os.File) error {
		return cb(json.NewDecoder(f))
	})
}
