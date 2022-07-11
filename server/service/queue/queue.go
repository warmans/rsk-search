package queue

import (
	gcloud "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	TaskImportCreateWorkspace           = "import:create_workspace"
	TaskImportCreateWav                 = "import:create_workspace"
	TaskImportMachineTranscribe         = "import:machine_transcribe"
	TaskImportChunkMachineTranscription = "import:chunk_machine_transcription"
	TaskImportSplitChunks               = "import:split_mp3_chunks"
	TaskImportPublish                   = "import:publish"
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
	fs     afero.Fs
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

	var p *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	localMp3Path := fmt.Sprintf("%s/%s.mp3", p.WorkingDir(), p.EpID)

	fileAlreadyExists, err := afero.Exists(q.fs, localMp3Path)
	if err != nil {
		return err
	}
	if fileAlreadyExists {
		// for some reason this job was retried but the files are already there.
		return q.DispatchCreateWav(ctx, p)
	}
	if err := q.fs.Mkdir(p.WorkingDir(), os.ModeDir); err != nil {
		return err
	}

	q.logger.Info("Fetching MP3 to local working dir", zap.String("mp3", p.Mp3URI))
	localMP3, err := q.fs.Create(localMp3Path)
	if err != nil {
		return err
	}
	defer localMP3.Close()

	resp, err := http.DefaultClient.Get(p.Mp3URI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(localMP3, resp.Body); err != nil {
		return errors.Wrap(err, "failed to fetch mp3")
	}

	q.TryUpdateImportLog(ctx, p.ID, t.Type(), "Workspace + local MP3 created: %s", localMp3Path)

	// move to next stage
	return q.DispatchCreateWav(ctx, p)
}

func (q *ImportQueue) HandleCreateWav(ctx context.Context, t *asynq.Task) error {
	var p *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	localMp3Path := fmt.Sprintf("%s/%s.mp3", p.WorkingDir(), p.EpID)
	localWavPath := fmt.Sprintf("%s/%s.wav", p.WorkingDir(), p.EpID)

	cmd := exec.Command("ffmpeg", fmt.Sprintf(`-i "%s"`, localMp3Path), `-ac 1`, localWavPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "failed to exec ffmpeg")
		q.logger.Error(err.Error(), zap.Error(err), zap.Strings("output", strings.Split(fmt.Sprint(out), "\n")))
		return err
	}

	q.TryUpdateImportLog(ctx, p.ID, t.Type(), "WAV file created: %s", localWavPath)

	return nil // dispatch machine translate stage
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

// TryUpdateImportLog is a best effort to get import info back into the DB. But if it fails, just log the error instead.
func (q *ImportQueue) TryUpdateImportLog(ctx context.Context, id string, stage string, format string, params ...interface{}) {
	if err := q.rw.WithStore(func(s *rw.Store) error {
		return s.PushTscriptImportLog(
			ctx,
			id,
			&models.TscriptImportLog{
				Stage: stage,
				Msg:   fmt.Sprintf(format, params...),
			},
		)
	}); err != nil {
		q.logger.Error(
			"Failed to update import log",
			zap.String("id", id),
			zap.String("stage", string(stage)),
			zap.String("msg", fmt.Sprintf(format, params...)),
		)
	}
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
