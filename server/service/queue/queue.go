package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/assemblyai"
	"github.com/warmans/rsk-search/pkg/flag"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/speech2text"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"path"
	"time"
)

const (
	TaskImportCreateWorkspace   = "import:create_workspace"
	TaskImportMachineTranscribe = "import:machine_transcribe"
	TaskImportSplitChunks       = "import:split_mp3_chunks"
	TaskImportPublish           = "import:publish"
)

type ImportQueueConfig struct {
	Addr          string
	WorkingDir    string
	KeepFiles     bool
	AssemblyAIKey string
}

func (c *ImportQueueConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Addr, prefix, "import-redis-addr", "localhost:6379", "redis address to use for queue backend")
	flag.StringVarEnv(fs, &c.WorkingDir, prefix, "import-work-dir", "./var/imports", "location to store in progress import artifacts")
	flag.BoolVarEnv(fs, &c.KeepFiles, prefix, "import-keep-files", false, "do not remove files after publish")
	flag.StringVarEnv(fs, &c.AssemblyAIKey, prefix, "import-assembly-ai-key", "", "API key for assemblyAI")
}

type ImportPipeline interface {
	StartNewImport(ctx context.Context, tscriptImport *models.TscriptImport) error
}

func NewImportQueue(
	logger *zap.Logger,
	filesystem afero.Fs,
	rw *rw.Conn,
	assemblyAi *assemblyai.Client,
	cfg *ImportQueueConfig,
	mediaBasePath string) *ImportQueue {
	a := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{
			Concurrency: 3,
			Logger:      asynqZapLogger(logger),
		},
	)
	return &ImportQueue{
		logger:        logger,
		cfg:           cfg,
		srv:           a,
		client:        asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Addr}),
		rw:            rw,
		fs:            filesystem,
		assemblyAi:    assemblyAi,
		mediaBasePath: mediaBasePath,
	}
}

type ImportQueue struct {
	logger     *zap.Logger
	cfg        *ImportQueueConfig
	srv        *asynq.Server
	client     *asynq.Client
	rw         *rw.Conn
	assemblyAi *assemblyai.Client
	fs         afero.Fs
	// main media directory used for media serving from the HTTP server
	// we will read from here and write back to the chunks dir
	mediaBasePath string
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
	q.logger.Info("Enqueue create workspace...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWorkspace, payload), asynq.Timeout(time.Minute*10), asynq.MaxRetry(5))
	return err
}

func (q *ImportQueue) DispatchMachineTranscribe(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Info("Enqueue machine transcribe...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportMachineTranscribe, payload, asynq.Timeout(time.Hour*3), asynq.MaxRetry(-1)))
	return err
}

func (q *ImportQueue) DispatchSplitAudioChunks(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Info("Enqueue split audio chunks...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportSplitChunks, payload), asynq.Timeout(time.Hour*1), asynq.MaxRetry(5))
	return err
}

func (q *ImportQueue) DispatchPublish(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Info("Enqueue publish...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportPublish, payload), asynq.Timeout(time.Minute*2), asynq.MaxRetry(2))
	return err
}

func (q *ImportQueue) HandleCreateWorkspace(ctx context.Context, t *asynq.Task) error {

	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	localMp3Path := path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.Mp3())

	// remove if it already exists and start again
	err := q.fs.RemoveAll(tsImport.WorkingDir(q.cfg.WorkingDir))
	if err != nil {
		return err
	}
	if err := q.fs.Mkdir(tsImport.WorkingDir(q.cfg.WorkingDir), 0755); err != nil {
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

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Workspace created: %s", q.cfg.WorkingDir)

	return q.DispatchMachineTranscribe(ctx, tsImport)
}

func (q *ImportQueue) HandleCreateMachineTranscription(ctx context.Context, t *asynq.Task) error {

	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	rawOutputPath := path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.MachineTranscript())
	if err := q.fs.RemoveAll(rawOutputPath); err != nil {
		return err
	}
	outputFile, err := q.fs.Create(rawOutputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	q.logger.Info("Starting speech 2 text...")
	resp, err := q.assemblyAi.Transcribe(ctx, &assemblyai.TranscribeRequest{AudioURL: tsImport.Mp3URI, SpeakerLabels: true})
	if err != nil {
		q.logger.Error("Failed assemblyai text", zap.Error(err))
		return err
	}
	if err := assemblyai.ToFlatFile(resp, outputFile); err != nil {
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Machine transcription completed: %s", tsImport.WAV())

	if err := outputFile.Sync(); err != nil {
		return err
	}
	if _, err := outputFile.Seek(0, 0); err != nil {
		return err
	}

	q.logger.Info("Starting chunk splitting...")
	chunkedOutputPath := path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.ChunkedMachineTranscript())
	if err := q.fs.RemoveAll(chunkedOutputPath); err != nil {
		return err
	}
	outFile, err := q.fs.Create(chunkedOutputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := speech2text.MapChunksFromRawTranscript(tsImport.EpID, tsImport.EpName, outputFile, outFile); err != nil {
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Chunks created: %s", tsImport.ChunkedMachineTranscript())

	return q.DispatchPublish(ctx, tsImport)
}

func (q *ImportQueue) HandlePublish(ctx context.Context, t *asynq.Task) error {
	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	q.logger.Info("Import chunks to DB")
	tscript := &models.ChunkedTranscript{}
	if err := util.WithReadJSONFileDecoder(path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.ChunkedMachineTranscript()), func(dec *json.Decoder) error {
		return dec.Decode(tscript)
	}); err != nil {
		q.logger.Error("failed to read JSON file", zap.Error(err))
		return err
	}
	if err := q.rw.WithStore(func(s *rw.Store) error {
		if err := s.InsertOrIgnoreTscript(context.Background(), tscript); err != nil {
			q.logger.Error("failed insert tscript", zap.Error(err))
			return err
		}
		if err := s.CompleteTscriptImport(ctx, tsImport.ID); err != nil {
			q.logger.Error("failed to complete tscript import", zap.Error(err))
			return err
		}
		if q.cfg.KeepFiles {
			return nil
		}
		return q.fs.RemoveAll(tsImport.WorkingDir(q.cfg.WorkingDir))
	}); err != nil {
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Publish complete")
	return nil
}

func (q *ImportQueue) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskImportCreateWorkspace, q.HandleCreateWorkspace)
	mux.HandleFunc(TaskImportMachineTranscribe, q.HandleCreateMachineTranscription)
	mux.HandleFunc(TaskImportPublish, q.HandlePublish)
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

func (q *ImportQueue) copyLocalFile(srcPath string, destPath string) error {
	localChunk, err := q.fs.Open(srcPath)
	if err != nil {
		return err
	}
	defer localChunk.Close()

	destFile, err := q.fs.Create(destPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(destFile, localChunk); err != nil {
		return err
	}

	return destFile.Close()
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
