package mediacache

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"syscall"
)

type Config struct {
	DataPath string
	Disabled bool
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.DataPath, prefix, "media-cache-data-dir", "./var/cache/media", "Directory to store cached media")
	flag.BoolVarEnv(fs, &c.Disabled, prefix, "media-cache-disabled", false, "Do no cache files")
}

type Cache struct {
	cfg    Config
	logger *zap.Logger
}

func NewCache(cfg Config, log *zap.Logger) (*Cache, error) {
	if cfg.DataPath == "" {
		return nil, errors.New("mediacache: must specify --data-path")
	}
	return &Cache{cfg: cfg, logger: log.With(zap.String("component", "media_cache"))}, nil
}

func (c *Cache) Get(key string, writeTo io.Writer, fetchFn func(writer io.Writer) error) (bool, error) {
	if c.cfg.Disabled {
		return false, fetchFn(writeTo)
	}
	filePath := path.Join(c.cfg.DataPath, key)
	f, err := os.Open(filePath)
	if err == nil {
		defer f.Close()
		if _, err = io.Copy(writeTo, f); err == nil {
			return true, nil
		}
		c.logger.Error("failed to write to writer", zap.Error(err))
		return false, fetchFn(writeTo)
	}
	if !errors.Is(err, os.ErrNotExist) {
		c.logger.Error("failed to open cached file", zap.String("file_path", filePath), zap.Error(err))
		return false, fetchFn(writeTo)
	}

	// cached file doesn't exist
	err = func() error {
		newFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			c.logger.Error("failed to create cached file", zap.String("file_path", filePath), zap.Error(err))
			return fetchFn(writeTo)
		}
		defer func() {
			if err := newFile.Close(); err != nil {
				panic(fmt.Sprintf("failed to close file after write: %s", err.Error()))
			}
		}()
		if err = syscall.Flock(int(newFile.Fd()), syscall.LOCK_EX); err != nil {
			c.logger.Error("failed to lock file for writing", zap.String("file_path", filePath), zap.Error(err))
			return fetchFn(writeTo)
		}
		defer func() {
			if err := syscall.Flock(int(newFile.Fd()), syscall.LOCK_UN); err != nil {
				panic(fmt.Sprintf("failed to unlock file after write: %s", err.Error()))
			}
		}()
		err = fetchFn(io.MultiWriter(writeTo, newFile))
		if err != nil {
			if rmErr := os.Remove(filePath); rmErr != nil {
				c.logger.Error("failed to remove cached file after write error", zap.String("file_path", filePath), zap.Error(rmErr))
			}
		}
		return err
	}()
	if err != nil {
		if err := os.Remove(filePath); err != nil {
			c.logger.Error("failed to remove cached file after write error", zap.String("file_path", filePath), zap.Error(err))
		}
	}
	return false, err
}
