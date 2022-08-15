package coffee

import (
	"context"
	"fmt"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"time"
)

func NewWorker(client *Client, db *rw.Conn, logger *zap.Logger, cfg *Config) *Worker {
	return &Worker{
		db:     db,
		stop:   make(chan struct{}),
		logger: logger.With(zap.String("component", "supporter worker")),
		cfg:    cfg,
		client: client,
	}
}

// Worker is not concurrency safe. Would need a global lock or similar if running multiple instances.
type Worker struct {
	db       *rw.Conn
	stop     chan struct{}
	stopping bool
	logger   *zap.Logger
	cfg      *Config
	client   *Client
}

func (w *Worker) Start() error {
	ticker := time.NewTicker(time.Second * time.Duration(w.cfg.SupporterSyncInterval))
	defer ticker.Stop()

	w.logger.Info("Starting supporter sync worker...")
	for {
		select {
		case <-w.stop:
			return nil
		default:
		}
		select {
		case <-ticker.C:
			w.logger.Info("Fetching supporters...")
			supporters, err := w.client.Supporters()
			if err != nil {
				w.logger.Error("Failed sync supporters from coffee", zap.Error(err))
			}
			supNames := []string{}
			for _, v := range supporters.Data {
				supNames = append(supNames, v.Name)
			}
			w.logger.Info("Found supporters", zap.Strings("names", supNames))
			err = w.db.WithStore(func(s *rw.Store) error {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				return s.SetAuthorsAsSupporters(ctx, supNames)
			})
			if err != nil {
				w.logger.Error("Failed to update local DB with supporters", zap.Error(err))
			}
		case <-w.stop:
			return nil
		}
	}
}

func (w *Worker) Stop(ctx context.Context) error {
	w.logger.Info("Stopping reward worker...")

	w.stopping = true
	stopped := make(chan struct{})
	go func() {
		close(w.stop)
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout stopping reward worker")
	case <-stopped:
		return nil
	}
}
