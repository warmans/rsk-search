package util

import (
	"encoding/json"
	"fmt"
	"os"
)

func WithCreateOrReplaceFile(path string, cb func(f *os.File) error) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WithCreateJSONFileEncoder(path string, cb func(enc *json.Encoder) error) error {
	return WithNewFile(path, func(f *os.File) error {
		enc := json.NewEncoder(f)
		enc.SetIndent("  ", "  ")
		return cb(enc)
	})
}

func WithReplaceJSONFileEncoder(path string, cb func(enc *json.Encoder) error) error {
	return WithExistingFile(path, func(f *os.File) error {
		if err := f.Truncate(0); err != nil {
			return err
		}
		if _, err := f.Seek(0, 0); err != nil {
			return err
		}
		enc := json.NewEncoder(f)
		enc.SetIndent("  ", "  ")
		return cb(enc)
	})
}

func WithReadJSONFileDecoder(path string, cb func(dec *json.Decoder) error) error {
	return WithExistingFile(path, func(f *os.File) error {
		return cb(json.NewDecoder(f))
	})
}
