package queue

import (
	gcloud "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
)

const (
	TaskImportCreateWorkspace = "import:create_workspace"
	TaskImportCreateWav       = "import:" + string(models.TscriptImportStageCreateWAV)
)

type ImportDispatcher interface {
	DispatchCreateWav(ctx context.Context, importID string) error
}

func NewImportQueue(logger *zap.Logger, redisAddr string) *ImportQueue {
	a := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 3,
			Logger:      asynqZapLogger(logger),
		},
	)
	return &ImportQueue{
		logger: logger,
		srv:    a,
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
	}
}

type ImportQueue struct {
	logger *zap.Logger
	srv    *asynq.Server
	client *asynq.Client
	rw     *rw.Conn
	gcloud *gcloud.Client
}

func (q *ImportQueue) DispatchCreateWorkspace(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWorkspace, payload))
	return err
}

func (q *ImportQueue) DispatchCreateWav(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWav, payload))
	return err
}

func (q *ImportQueue) HandleCreateWorkspace(ctx context.Context, t *asynq.Task) error {
	//todo: create bucket in gcloud for import files
	return nil
}

func (q *ImportQueue) HandleCreateWav(ctx context.Context, t *asynq.Task) error {
	var p *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	q.logger.Info("Converting .mp3 to .wav", zap.String("mp3", p.Mp3URI))

	//todo: create the wav and store it back to google cloud

	return q.rw.WithStore(func(s *rw.Store) error {
		return s.SetTscriptImportStage(
			ctx,
			p.ID,
			models.TscriptImportStageMachineTranscribe,
			&models.TscriptImportLog{
				Stage: models.TscriptImportStageCreateWAV,
				Msg:   fmt.Sprintf("Success"),
				Data: map[string]interface{}{
					"path": "some wav path",
				},
			},
		)
	})
}

func (q *ImportQueue) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskImportCreateWorkspace, q.HandleCreateWorkspace)
	mux.HandleFunc(TaskImportCreateWav, q.HandleCreateWav)
	return q.srv.Start(mux)
}

func (q *ImportQueue) Stop() {
	q.srv.Stop()
}

func asynqZapLogger(zap *zap.Logger) asynq.Logger {
	return &ZapLogger{zap: zap}
}

type ZapLogger struct {
	zap *zap.Logger
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.zap.Debug(fmt.Sprint(args...))
}

func (z ZapLogger) Info(args ...interface{}) {
	z.zap.Info(fmt.Sprint(args...))
}

func (z ZapLogger) Warn(args ...interface{}) {
	z.zap.Warn(fmt.Sprint(args...))
}

func (z ZapLogger) Error(args ...interface{}) {
	z.zap.Error(fmt.Sprint(args...))
}

func (z ZapLogger) Fatal(args ...interface{}) {
	z.zap.Fatal(fmt.Sprint(args...))
}
