package speech2text

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"go.uber.org/zap"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io"
	"os"
	"path"
)

type GcloudConfig struct {
	Bucket         string
	WavStoragePath string
}

func (c *GcloudConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Bucket, prefix, "speech2text-gcloud-bucket", "scrimpton-raw-audio", "bucket to store WAVs when transcribing")
	flag.StringVarEnv(fs, &c.WavStoragePath, prefix, "speech2text-gcloud-wav-store", "wav", "sub dir within the bucket")
}

func NewGcloud(
	logger *zap.Logger,
	sto *storage.Client,
	spe *speech.Client,
	cfg *GcloudConfig,
) *Gcloud {
	return &Gcloud{
		logger:         logger,
		storage:        sto,
		speech:         spe,
		bucket:         cfg.Bucket,
		wavStoragePath: cfg.WavStoragePath,
	}
}

type Gcloud struct {
	logger         *zap.Logger
	storage        *storage.Client
	speech         *speech.Client
	bucket         string
	wavStoragePath string
}

func (g *Gcloud) GenerateText(ctx context.Context, localWavPath string, outputWriter io.Writer) error {

	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return fmt.Errorf("no google application credentials in env")
	}

	gsUtilURL, err := g.uploadWav(ctx, localWavPath)
	if err != nil {
		return err
	}

	g.logger.Info("Uploaded WAV", zap.String("url", gsUtilURL))

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_LINEAR16,
			LanguageCode:               "en-US",
			AudioChannelCount:          1,
			EnableWordTimeOffsets:      true,
			EnableAutomaticPunctuation: true,
			Model:                      "video",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gsUtilURL},
		},
	}

	g.logger.Info("Starting LongRunningRecognize")
	op, err := g.speech.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}

	g.logger.Info("Waiting for operation to complete...")
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			for _, v := range alt.Words {
				if _, err := fmt.Fprintf(outputWriter, "#OFFSET: %d\n", v.StartTime.Seconds); err != nil {
					return err
				}
				break
			}
			if _, err := fmt.Fprintf(outputWriter, "Unknown: %s\n", alt.Transcript); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Gcloud) uploadWav(ctx context.Context, localWavPath string) (string, error) {

	localFile, err := os.Open(localWavPath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	// ctx timeout?

	writer := g.storage.Bucket(g.bucket).Object(path.Join(g.wavStoragePath, path.Base(localWavPath))).NewWriter(ctx)

	if _, err := io.Copy(writer, localFile); err != nil {
		writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("gs://%s/%s/%s", g.bucket, g.wavStoragePath, path.Base(localWavPath)), nil
}
