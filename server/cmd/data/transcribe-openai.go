package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func TranscribeOpenAICommand() *cobra.Command {

	var inputPath string
	var workspaceBaseDir string
	var outputPath string
	var mergeOnly bool

	cmd := &cobra.Command{
		Use:   "transcribe-openai",
		Short: "create a machine transcription using openai wisper-1",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			apiKey := os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				logger.Fatal("OPENAI_API_KEY was not in environment")
			}

			workspace := path.Join(workspaceBaseDir, strings.TrimSuffix(path.Base(inputPath), ".mp3"))

			if mergeOnly {
				return mergeDialog(workspace)
			}

			if fileInfo, err := os.Stat(workspace); errors.Is(err, os.ErrNotExist) {
				if err := os.Mkdir(workspace, 0755); err != nil {
					return fmt.Errorf("failed to create workspace: %w", err)
				}
			} else {
				if err != nil {
					return err
				}
				if !fileInfo.IsDir() {
					return fmt.Errorf("%s shold be a directory, but was not", workspace)
				}
			}
			logger.Info("Workspace setup OK", zap.String("workspace", workspace))

			logger.Info("Splitting file...", zap.String("workspace", workspace))
			if err := splitFile(inputPath, workspace); err != nil {
				return fmt.Errorf("failed to split file: %w", err)
			}
			logger.Info("Processing files...", zap.String("workspace", workspace))
			entries, err := os.ReadDir(workspace)
			if err != nil {
				return fmt.Errorf("failed to read dir: %w", err)
			}
			for _, v := range entries {
				if v.IsDir() || !strings.HasSuffix(v.Name(), ".mp3") {
					continue
				}
				if err = dispatchFile(logger, path.Join(workspace, v.Name()), apiKey); err != nil {
					return err
				}
				logger.Info("Completed file", zap.String("file", v.Name()))
			}
			logger.Info("all files processed, merging files")

			return mergeDialog(workspace)
		},
	}

	cmd.Flags().StringVarP(&inputPath, "input-raw-audio-file", "i", path.Join(os.Getenv("MEDIA_BASE_PATH"), "episode", "xfm-S2E01.mp3"), "Input audio file.")
	cmd.Flags().StringVarP(&workspaceBaseDir, "workspace", "w", path.Join(os.Getenv("MEDIA_BASE_PATH"), "workspace"), "path to stage files for processing")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "./var/", "Dump output to given path.")
	cmd.Flags().BoolVarP(&mergeOnly, "merge-only", "m", false, "Only merge existing files")

	return cmd
}

// ffmpeg -i somefile.mp3 -f segment -segment_time 3 -c copy out%03d.mp3
func splitFile(inputFile string, workspace string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-f", "segment",
		"-segment_time", "00:30:00",
		"-c", "copy",
		strings.TrimSuffix(path.Base(inputFile), ".mp3")+".%03d.mp3",
	)
	cmd.Dir = workspace
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to shell out to ffmpeg for file: %s", inputFile)
	}
	return nil
}

//func getAudioDurationMs(audioFilePath string) (int64, error) {
//	cmd := exec.Command("mp3info", "-p", "%S", audioFilePath)
//	out, err := cmd.CombinedOutput()
//	if err != nil {
//		return 0, errors.Wrapf(err, "failed to shell out to mp3info (is it installed?) for file: %s (raw: %s)", audioFilePath, string(out))
//	}
//	intVal, err := strconv.Atoi(string(out))
//	if err != nil {
//		return 0, errors.Wrapf(err, "failed to convert mp3info output to int: %s", string(out))
//	}
//	return int64(intVal) * 1000, nil
//}

func dispatchFile(logger *zap.Logger, inputFilePath string, openApiKey string) error {

	outputFilePath := fmt.Sprintf("%s.json", inputFilePath)
	if _, err := os.Stat(outputFilePath); err == nil {
		return nil // already done
	}

	logger.Info("Creating output file", zap.String("output-file", outputFilePath))
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			logger.Error("failed to close file ", zap.Error(err))
		}
	}(outputFile)

	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", inputFilePath, err)
	}
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			logger.Error("failed to close file ", zap.Error(err))
		}
	}(inputFile)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("model", "whisper-1"); err != nil {
		return fmt.Errorf("failed to add field: %w", err)
	}
	if err := writer.WriteField("timestamp_granularities[]", "segment"); err != nil {
		return fmt.Errorf("failed to add field: %w", err)
	}
	if err := writer.WriteField("response_format", "verbose_json"); err != nil {
		return fmt.Errorf("failed to add field: %w", err)
	}
	part, err := writer.CreateFormFile("file", path.Base(inputFilePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, inputFile); err != nil {
		return fmt.Errorf("failed to copy input file: %w", err)
	}
	if err = writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/audio/transcriptions", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", openApiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send api request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if _, err = io.Copy(outputFile, resp.Body); err != nil {
		return fmt.Errorf("failed to dump response: %w", err)
	}

	return nil
}

func mergeDialog(workspace string) error {

	files, err := os.ReadDir(workspace)
	if err != nil {
		return fmt.Errorf("failed to list workspace files: %w", err)
	}

	dialog := []models.Dialog{}
	var startPos int64
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") || f.Name() == "dialog.json" {
			continue
		}

		openAIData := struct {
			Segments []struct {
				Position int64   `json:"id"`
				Start    float64 `json:"start"`
				End      float64 `json:"end"`
				Text     string  `json:"text"`
			} `json:"segments"`
		}{}
		err = util.WithReadJSONFileDecoder(path.Join(workspace, f.Name()), func(dec *json.Decoder) error {
			return dec.Decode(&openAIData)
		})
		if err != nil {
			return fmt.Errorf("failed to decode file: %w", err)
		}
		if len(dialog) > 0 {
			startPos = dialog[len(dialog)-1].Position
		}
		for _, segment := range openAIData.Segments {
			if text := strings.ToLower(strings.Trim(strings.TrimSpace(segment.Text), ".")); text == "song" || text == "songs" {
				continue
			}
			dialog = append(
				dialog,
				models.Dialog{
					Position:  startPos + segment.Position,
					Timestamp: time.Duration(segment.Start*1000) * time.Millisecond,
					Content:   strings.TrimSpace(segment.Text),
					Duration:  time.Duration(segment.End*1000)*time.Millisecond - time.Duration(segment.Start*1000)*time.Millisecond,
				},
			)
		}
	}

	return util.WithReplaceJSONFileEncoder(path.Join(workspace, "dialog.json"), func(enc *json.Encoder) error {
		return enc.Encode(dialog)
	})
}
