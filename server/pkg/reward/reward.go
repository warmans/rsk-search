package reward

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"time"
)

// RewardSpacing - every N approved contributions a reward will be triggered
const RewardSpacing = 5

type Config struct {
	CheckInterval int64
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.Int64VarEnv(fs, &c.CheckInterval, prefix, "reward-check-interval-seconds", 60, "check for pending rewards every N seconds")
}

func NewWorker(db *rw.Conn, logger *zap.Logger, cfg Config) *Worker {
	return &Worker{
		db:     db,
		stop:   make(chan struct{}, 0),
		logger: logger.With(zap.String("component", "reward worker")),
		cfg:    cfg,
	}
}

// Worker is not concurrency safe. Would need a global lock or similar if running multiple instances.
type Worker struct {
	db       *rw.Conn
	stop     chan struct{}
	stopping bool
	logger   *zap.Logger
	cfg      Config
}

func (w *Worker) Start() error {
	ticker := time.NewTicker(time.Second * time.Duration(w.cfg.CheckInterval))
	defer ticker.Stop()

	w.logger.Info("Starting reward worker...")
	for {
		select {
		case <-w.stop:
			return nil
		default:
		}
		select {
		case <-ticker.C:
			w.logger.Debug("Calculating rewards")
			if err := w.calculateRewards(); err != nil {
				w.logger.Error("Failed to run rewards", zap.Error(err))
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

func (w *Worker) calculateRewards() error {

	var awardsRequired []*models.AuthorReward

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err := w.db.WithStore(func(s *rw.Store) error {
		var err error
		awardsRequired, err = s.ListRequiredAuthorRewards(ctx, RewardSpacing)
		if err != nil {
			return errors.Wrap(err, "failed to list required rewards")
		}
		w.logger.Debug("Num rewards required", zap.Int("num", len(awardsRequired)))
		for k, a := range awardsRequired {
			if awardsRequired[k].ID, err = s.CreatePendingReward(ctx, a.AuthorID, a.Threshold); err != nil {
				return errors.Wrap(err, "failed to create pending reward")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, a := range awardsRequired {
		w.logger.Info(
			"Created reward",
			zap.String("id", a.ID),
			zap.String("author_id", a.AuthorID),
			zap.Int32("threshold", a.Threshold),
		)
	}
	return nil
}
