package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/assemblyai"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func TranscribeAssemblyAICmd() *cobra.Command {

	var inputURL string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "transcribe-assembly-ai",
		Short: "create a machine transcription using Assembly AI API",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			apiKey := os.Getenv("ASSEMBLY_AI_KEY")
			if apiKey == "" {
				logger.Fatal("Assembly AI API key was not in environment: ASSEMBLY_AI_KEY")
			}

			logger.Info("Creating output file", zap.String("output-file", outputPath))
			outputFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			ctx, done := context.WithTimeout(context.Background(), time.Minute*30)
			defer done()

			client := assemblyai.NewClient(logger, http.DefaultClient, &assemblyai.Config{AccessToken: apiKey})
			result, err := client.Transcribe(ctx, &assemblyai.TranscribeRequest{AudioURL: inputURL, SpeakerLabels: true})
			if err != nil {
				return err
			}

			logger.Info("Completed!")

			// dump output to JSON
			enc := json.NewEncoder(outputFile)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		},
	}

	cmd.Flags().StringVarP(&inputURL, "input-audio-url", "i", "https://storage.googleapis.com/scrimpton-raw-audio/sample-S0E0.mp3", "Input audio file.")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "./var/aai-output-"+shortuuid.New()+".json", "Dump output to given path.")

	return cmd
}

func AssemblyAI2Dialog() *cobra.Command {

	var intputPath string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "convert-assembly-ai",
		Short: "create a machine transcription using Assembly AI API",
		RunE: func(cmd *cobra.Command, args []string) error {

			if intputPath == "" {
				return fmt.Errorf("input path not set")
			}

			inFile, err := os.Open(intputPath)
			if err != nil {
				return err
			}

			outFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			resp := &assemblyai.TranscriptionStatusResponse{}
			if err := json.NewDecoder(inFile).Decode(resp); err != nil {
				return err
			}

			return assemblyai.ToFlatFile(resp, outFile)
		},
	}

	cmd.Flags().StringVarP(&intputPath, "input-path", "i", "", "raw file created with transcribe-assembly-ai")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "./var/aai-output-"+shortuuid.New()+".tscript.txt", "Dump output to given path.")

	return cmd
}
