package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/speech2text"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os/exec"
	"path"
	"strings"
)

const (
	TaskImportCreateWorkspace           = "import:create_workspace"
	TaskImportCreateWav                 = "import:create_wav"
	TaskImportMachineTranscribe         = "import:machine_transcribe"
	TaskImportChunkMachineTranscription = "import:chunk_machine_transcription"
	TaskImportSplitChunks               = "import:split_mp3_chunks"
	TaskImportPublish                   = "import:publish"
)

type ImportqueueConfig struct {
	Addr         string
	WorkingDir   string
	GcloudBucket string
	GcloudWavDir string
}

func (c *ImportqueueConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Addr, prefix, "redis-addr", "localhost:6379", "redis address to use for queue backend")
	flag.StringVarEnv(fs, &c.WorkingDir, prefix, "work-dir", "./var/imports", "location to store in progress import artifacts")
	flag.StringVarEnv(fs, &c.GcloudBucket, prefix, "gcloud-bucket", "scrimpton-raw-audio", "bucket to upload data where needed")
	flag.StringVarEnv(fs, &c.GcloudWavDir, prefix, "gcloud-wav-dir", "wav", "bucket to upload data where needed")
}

type ImportPipeline interface {
	StartNewImport(ctx context.Context, tscriptImport *models.TscriptImport) error
}

func NewImportQueue(
	logger *zap.Logger,
	filesystem afero.Fs,
	rw *rw.Conn,
	speech2text *speech2text.Gcloud,
	cfg *ImportqueueConfig) *ImportQueue {
	a := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{
			Concurrency: 3,
			Logger:      asynqZapLogger(logger),
		},
	)
	return &ImportQueue{
		logger:      logger,
		srv:         a,
		client:      asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Addr}),
		rw:          rw,
		fs:          filesystem,
		workDirRoot: cfg.WorkingDir,
		speech2text: speech2text,
	}
}

type ImportQueue struct {
	logger      *zap.Logger
	srv         *asynq.Server
	client      *asynq.Client
	rw          *rw.Conn
	speech2text *speech2text.Gcloud
	fs          afero.Fs
	workDirRoot string
}

// StartNewImport implements a simple interface for the server. It doesn't need to know the steps to run an import.
func (q *ImportQueue) StartNewImport(ctx context.Context, tscriptImport *models.TscriptImport) error {
	return q.DispatchCreateWorkspace(ctx, tscriptImport)
}

func (q *ImportQueue) DispatchCreateWorkspace(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue create workspace...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWorkspace, payload))
	return err
}

func (q *ImportQueue) DispatchCreateWav(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue create WAV...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWav, payload))
	return err
}

func (q *ImportQueue) DispatchMachineTranscribe(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue machine transcribe...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportMachineTranscribe, payload))
	return err
}

func (q *ImportQueue) HandleCreateWorkspace(ctx context.Context, t *asynq.Task) error {

	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	localMp3Path := fmt.Sprintf("%s/%s.mp3", tsImport.WorkingDir(q.workDirRoot), tsImport.EpID)

	fileAlreadyExists, err := afero.Exists(q.fs, localMp3Path)
	if err != nil {
		return err
	}
	if fileAlreadyExists {
		// for some reason this job was retried but the files are already there.
		return q.DispatchCreateWav(ctx, tsImport)
	}
	if err := q.fs.Mkdir(tsImport.WorkingDir(q.workDirRoot), 0755); err != nil {
		return err
	}

	q.logger.Info("Creating local MP3", zap.String("path", localMp3Path))
	localMP3, err := q.fs.Create(localMp3Path)
	if err != nil {
		return err
	}
	defer localMP3.Close()

	q.logger.Info("Fetching MP3 to local working dir", zap.String("uri", tsImport.Mp3URI))
	resp, err := http.DefaultClient.Get(tsImport.Mp3URI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	q.logger.Info("Copying to file")
	if _, err := io.Copy(localMP3, resp.Body); err != nil {
		return errors.Wrap(err, "failed to fetch mp3")
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Workspace + local MP3 created: %s", localMp3Path)

	// move to next stage
	q.logger.Info("Dispatch create WAV task")
	return q.DispatchCreateWav(ctx, tsImport)
}

func (q *ImportQueue) HandleCreateWav(ctx context.Context, t *asynq.Task) error {
	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// remove wav if it was already created
	if err := q.fs.RemoveAll(path.Join(tsImport.WorkingDir(q.workDirRoot), tsImport.WAV())); err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", tsImport.Mp3(), `-ac`, `1`, tsImport.WAV())
	cmd.Dir = tsImport.WorkingDir(q.workDirRoot)

	q.logger.Info("Shelling out to ffmpeg to complete wav conversation")
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "failed to exec ffmpeg")
		q.logger.Error(err.Error(), zap.Error(err), zap.Strings("output", strings.Split(string(out), "\n")))
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "WAV file created: %s", tsImport.WAV())

	return nil // dispatch machine translate stage
}

func (q *ImportQueue) HandleCreateMachineTranscription(ctx context.Context, t *asynq.Task) error {

	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	f, err := q.fs.Create(path.Join(tsImport.WorkingDir(q.workDirRoot), tsImport.MachineTranscript()))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := q.speech2text.GenerateText(ctx, path.Join(tsImport.WorkingDir(q.workDirRoot), tsImport.WAV()), f); err != nil {
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Machine transcription completed: %s", tsImport.WAV())

	return nil
}

func (q *ImportQueue) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskImportCreateWorkspace, q.HandleCreateWorkspace)
	mux.HandleFunc(TaskImportCreateWav, q.HandleCreateWav)
	mux.HandleFunc(TaskImportMachineTranscribe, q.HandleCreateMachineTranscription)
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
			zap.String("stage", stage),
			zap.String("msg", fmt.Sprintf(format, params...)),
			zap.Error(err),
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
