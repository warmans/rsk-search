package queue

import (
	"cloud.google.com/go/storage"
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
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	TaskImportCreateWorkspace   = "import:create_workspace"
	TaskImportCreateWav         = "import:create_wav"
	TaskImportMachineTranscribe = "import:machine_transcribe"
	TaskImportSplitChunks       = "import:split_mp3_chunks"
	TaskImportPublish           = "import:publish"
)

type ImportQueueConfig struct {
	Addr                string
	WorkingDir          string
	GcloudAudioBucket   string
	GcloudChunkBucket   string
	GcloudWavDir        string
	PythonScriptDir     string
	KeepFiles           bool
	AssemblyAIKey       string
	Speech2TextStrategy string
}

func (c *ImportQueueConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Addr, prefix, "import-redis-addr", "localhost:6379", "redis address to use for queue backend")
	flag.StringVarEnv(fs, &c.WorkingDir, prefix, "import-work-dir", "./var/imports", "location to store in progress import artifacts")
	flag.StringVarEnv(fs, &c.GcloudAudioBucket, prefix, "import-gcloud-audio-bucket", "scrimpton-raw-audio", "bucket to upload data where needed")
	flag.StringVarEnv(fs, &c.GcloudChunkBucket, prefix, "import-gcloud-chunk-bucket", "scrimpton-chunked-audio", "bucket for storing audio chunks")
	flag.StringVarEnv(fs, &c.GcloudWavDir, prefix, "import-gcloud-wav-dir", "wav", "bucket to upload data where needed")
	flag.StringVarEnv(fs, &c.PythonScriptDir, prefix, "import-python-script-dir", "./script/audio-splitter", "location of python scripts used for audio splitting")
	flag.BoolVarEnv(fs, &c.KeepFiles, prefix, "import-keep-files", false, "do not remove files after publish")
	flag.StringVarEnv(fs, &c.AssemblyAIKey, prefix, "import-assembly-ai-key", "", "API key for assemblyAI")
	flag.StringVarEnv(fs, &c.Speech2TextStrategy, prefix, "import-speech-2-text-strategy", "assemblyai", "gclound or assemblyai")
}

type ImportPipeline interface {
	StartNewImport(ctx context.Context, tscriptImport *models.TscriptImport) error
}

func NewImportQueue(
	logger *zap.Logger,
	filesystem afero.Fs,
	rw *rw.Conn,
	speech2text *speech2text.Gcloud,
	assemblyAi *assemblyai.Client,
	gcloud *storage.Client,
	cfg *ImportQueueConfig) *ImportQueue {
	a := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Addr},
		asynq.Config{
			Concurrency: 3,
			Logger:      asynqZapLogger(logger),
		},
	)
	return &ImportQueue{
		logger:      logger,
		cfg:         cfg,
		srv:         a,
		client:      asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Addr}),
		rw:          rw,
		fs:          filesystem,
		gcloud:      gcloud,
		speech2text: speech2text,
		assemblyAi:  assemblyAi,
	}
}

type ImportQueue struct {
	logger      *zap.Logger
	cfg         *ImportQueueConfig
	srv         *asynq.Server
	client      *asynq.Client
	rw          *rw.Conn
	speech2text *speech2text.Gcloud
	assemblyAi  *assemblyai.Client
	gcloud      *storage.Client
	fs          afero.Fs
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
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWorkspace, payload), asynq.Timeout(time.Minute*10), asynq.MaxRetry(5))
	return err
}

func (q *ImportQueue) DispatchCreateWav(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue create WAV...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportCreateWav, payload), asynq.Timeout(time.Hour), asynq.MaxRetry(5))
	return err
}

func (q *ImportQueue) DispatchMachineTranscribe(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue machine transcribe...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportMachineTranscribe, payload, asynq.Timeout(time.Hour*3), asynq.MaxRetry(-1)))
	return err
}

func (q *ImportQueue) DispatchSplitAudioChunks(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue split audio chunks...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportSplitChunks, payload), asynq.Timeout(time.Hour*1), asynq.MaxRetry(5))
	return err
}

func (q *ImportQueue) DispatchPublish(ctx context.Context, tscriptImport *models.TscriptImport) error {
	payload, err := json.Marshal(tscriptImport)
	if err != nil {
		return err
	}
	q.logger.Debug("Enqueue publish...")
	_, err = q.client.EnqueueContext(ctx, asynq.NewTask(TaskImportPublish, payload))
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
	// move to next stage

	if q.cfg.Speech2TextStrategy == "assemblyai" {
		// wav is only needed for gcloud
		q.logger.Info("Dispatch machine transcribe task (skip WAV)")
		return q.DispatchMachineTranscribe(ctx, tsImport)
	} else {
		q.logger.Info("Dispatch create WAV task")
		return q.DispatchCreateWav(ctx, tsImport)
	}
}

func (q *ImportQueue) HandleCreateWav(ctx context.Context, t *asynq.Task) error {

	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// remove wav if it was already created
	if err := q.fs.RemoveAll(path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.WAV())); err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg", "-i", tsImport.Mp3(), `-ac`, `1`, tsImport.WAV())
	cmd.Dir = tsImport.WorkingDir(q.cfg.WorkingDir)

	q.logger.Info("Shelling out to ffmpeg to complete wav conversation")
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "failed to exec ffmpeg")
		q.logger.Error(err.Error(), zap.Error(err), zap.Strings("output", strings.Split(string(out), "\n")))
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "WAV file created: %s", tsImport.WAV())

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

	q.logger.Info("Starting speech 2 text...", zap.String("strategy", q.cfg.Speech2TextStrategy))
	switch q.cfg.Speech2TextStrategy {
	case "gcloud":
		if err := q.speech2text.GenerateText(ctx, path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.WAV()), outputFile); err != nil {
			q.logger.Error("Failed gcloud speech 2 text", zap.Error(err))
			return err
		}
	case "assemblyai":
		resp, err := q.assemblyAi.Transcribe(ctx, &assemblyai.TranscribeRequest{AudioURL: tsImport.Mp3URI, SpeakerLabels: true})
		if err != nil {
			q.logger.Error("Failed assemblyai text", zap.Error(err))
			return err
		}
		if err := assemblyai.ToFlatFile(resp, outputFile); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown speech 2 text strategy: %s", q.cfg.Speech2TextStrategy)
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

	if err := speech2text.MapChunksFromGoogleTranscript(tsImport.EpID, tsImport.EpName, outputFile, outFile); err != nil {
		return err
	}

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Chunks created: %s", tsImport.ChunkedMachineTranscript())

	return q.DispatchSplitAudioChunks(ctx, tsImport)
}

func (q *ImportQueue) HandleSplitAudioChunks(ctx context.Context, t *asynq.Task) error {
	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	chunkOutputDir := path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), "chunks")
	if err := q.fs.RemoveAll(chunkOutputDir); err != nil {
		return err
	}
	if err := q.fs.Mkdir(chunkOutputDir, 0755); err != nil {
		return err
	}

	chunkMetadataPath := path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.ChunkedMachineTranscript())

	var tscript *models.Tscript
	if err := util.WithReadJSONFileDecoder(chunkMetadataPath, func(dec *json.Decoder) error {
		if err := dec.Decode(&tscript); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	for _, chunk := range tscript.Chunks {

		args := []string{
			"-i", path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.Mp3()),
			"-vcodec", "copy",
			"-acodec", "copy",
			"-ss", fmt.Sprintf("%d", chunk.StartSecond),
		}
		if chunk.EndSecond > -1 {
			args = append(args, "-to", fmt.Sprintf("%d", chunk.EndSecond))
		}
		// else it will default to eof

		cmd := exec.Command(
			"ffmpeg",
			append(args, path.Join(chunkOutputDir, fmt.Sprintf("%s.mp3", chunk.ID)))...,
		)

		q.logger.Info("Shelling out to ffmpeg to extract Mp3 chunk", zap.String("chunk_id", chunk.ID), zap.Int64("strat", chunk.StartSecond))
		out, err := cmd.CombinedOutput()
		if err != nil {
			err = errors.Wrap(err, "failed to exec split-ep.py")
			q.logger.Error(err.Error(), zap.Error(err), zap.Strings("output", strings.Split(string(out), "\n")))
			return err
		}
	}

	//wavPath := path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.Mp3())
	//
	//cmd := exec.Command(
	//	"python3",
	//	path.Join(q.cfg.PythonScriptDir, "split-ep.py"),
	//	"--meta", chunkMetadataPath,
	//	"--outpath", chunkOutputDir,
	//	"--audio", wavPath,
	//)

	q.TryUpdateImportLog(ctx, tsImport.ID, t.Type(), "Mp3 split complete: %s", chunkOutputDir)

	return q.DispatchPublish(ctx, tsImport)
}

func (q *ImportQueue) HandlePublish(ctx context.Context, t *asynq.Task) error {
	var tsImport *models.TscriptImport
	if err := json.Unmarshal(t.Payload(), &tsImport); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// upload chunks
	files, err := afero.ReadDir(q.fs, path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), "chunks"))
	if err != nil {
		return err
	}
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		if !strings.HasSuffix(v.Name(), ".mp3") {
			continue
		}
		if err := q.copyLocalFileToBucket(ctx, q.cfg.GcloudChunkBucket, v.Name(), path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), "chunks", v.Name())); err != nil {
			return err
		}
		q.logger.Info("Copied chunk to bucket", zap.String("file", v.Name()), zap.String("bucket", q.cfg.GcloudChunkBucket))
	}

	q.logger.Info("Import chunks to DB")
	tscript := &models.Tscript{}
	if err := util.WithReadJSONFileDecoder(path.Join(tsImport.WorkingDir(q.cfg.WorkingDir), tsImport.ChunkedMachineTranscript()), func(dec *json.Decoder) error {
		return dec.Decode(tscript)
	}); err != nil {
		return err
	}
	if err := q.rw.WithStore(func(s *rw.Store) error {
		if err := s.InsertOrIgnoreTscript(context.Background(), tscript); err != nil {
			return err
		}
		if err := s.CompleteTscriptImport(ctx, tsImport.ID); err != nil {
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
	mux.HandleFunc(TaskImportCreateWav, q.HandleCreateWav)
	mux.HandleFunc(TaskImportMachineTranscribe, q.HandleCreateMachineTranscription)
	mux.HandleFunc(TaskImportSplitChunks, q.HandleSplitAudioChunks)
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

func (q *ImportQueue) copyLocalFileToBucket(ctx context.Context, destBucket string, destName string, srcPath string) error {
	localChunk, err := q.fs.Open(srcPath)
	if err != nil {
		return err
	}
	defer localChunk.Close()

	remoteChunk := q.gcloud.Bucket(q.cfg.GcloudChunkBucket).Object(destName).NewWriter(ctx)
	defer remoteChunk.Close()

	if _, err := io.Copy(remoteChunk, localChunk); err != nil {
		return err
	}
	return nil
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
