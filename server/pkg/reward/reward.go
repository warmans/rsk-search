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

// it's a bit tricky to make this configurable as adjusting it could cause weird unexpected reward behavior.
const rewardThreshold = 5

const RewarderNoop = "noop"
const RewarderRedditGold = "reddit-gold"

type Cfg struct {
	CheckInterval int64
	Rewarder      string
}

func (c *Cfg) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.Int64VarEnv(fs, &c.CheckInterval, prefix, "reward-check-interval-seconds", 60, "check for pending rewards every N seconds")
	flag.StringVarEnv(fs, &c.Rewarder, prefix, "reward-strategy", RewarderNoop, "implementation of rewarder (noop, reddit-gold)")
}

func NewWorker(db *rw.Conn, logger *zap.Logger, cfg Cfg, rewarder Rewarder) *Worker {
	return &Worker{
		db:       db,
		stop:     make(chan struct{}, 0),
		logger:   logger.With(zap.String("component", "reward worker")),
		cfg:      cfg,
		rewarder: rewarder,
	}
}

// Worker is not concurrency safe. Would need a global lock or similar if running multiple instances.
type Worker struct {
	db       *rw.Conn
	stop     chan struct{}
	stopping bool
	logger   *zap.Logger
	cfg      Cfg
	rewarder Rewarder
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
		awardsRequired, err = s.ListRequiredAuthorRewards(ctx, rewardThreshold)
		if err != nil {
			return err
		}
		w.logger.Debug("Num rewards required", zap.Int("num", len(awardsRequired)))
		for k, a := range awardsRequired {
			if awardsRequired[k].ID, err = s.CreatePendingReward(ctx, a.AuthorID, a.Threshold); err != nil {
				return err
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


// deprecated: make user claim reward instead.
func (w *Worker) giveRewards() error {

	var pendingRewards []*models.AuthorReward
	err := w.db.WithStore(func(s *rw.Store) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var err error
		if pendingRewards, err = s.ListPendingRewards(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	w.logger.Debug("Have pending rewards", zap.Int("num", len(pendingRewards)))
	for _, r := range pendingRewards {
		if w.stopping {
			return nil
		}
		err := w.db.WithStore(func(s *rw.Store) error {

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel()

			author, err := s.GetAuthor(ctx, r.AuthorID)
			if err != nil {
				return err
			}

			ident, err := author.DecodeIdentity()
			if err != nil {
				return errors.Wrap(err, "failed to decode identity")
			}

			rewardErr := w.rewarder.Reward(ident.ID)

			// at this point we're fucked if the DB becomes unreachable so try it a few times

			retry := 10
			for retry > 0 {
				retry--
				if rewardErr != nil {
					w.logger.Error("failed to issue reward", zap.Error(rewardErr))
					if err := s.FailReward(ctx, r.ID, rewardErr.Error()); err != nil {
						w.logger.Error("failed to fail reward!", zap.Error(err), zap.String("original_error", rewardErr.Error()), zap.String("id", r.ID))
						time.Sleep(time.Second * 5)
						continue
					}
					return nil
				}
				if err := s.ClaimReward(ctx, r.ID); err != nil {
					w.logger.Error("failed to confirm reward!", zap.Error(err), zap.String("id", r.ID))
					time.Sleep(time.Second * 5)
					continue
				}
				return nil
			}
			return fmt.Errorf("exhausted retries attempting to confirm or fail reward: %s", r.ID)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
