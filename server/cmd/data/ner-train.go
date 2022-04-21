package data

import (
	"bufio"
	"encoding/json"
	"github.com/jdkato/prose/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func NERTrainCmd() *cobra.Command {

	var inputFile string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "ner-train",
		Short: "Extract some raw training data that must be manually updated to train the model",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					panic("failed to sync logger: "+err.Error())
				}
			}()

			trainingData, err := os.Open(inputFile)
			if err != nil {
				return nil
			}
			defer trainingData.Close()

			dataRows := []prose.EntityContext{}

			scanner := bufio.NewScanner(trainingData)
			for scanner.Scan() {

				row := prose.EntityContext{}
				if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
					return err
				}
				dataRows = append(dataRows, row)
			}

			model := prose.ModelFromData("RSK", prose.UsingEntities(dataRows))
			return model.Write(outputFile)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "./var/data/ner/training-data.json", "Input file")
	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "./var/data/ner/rsk-model", "Output file")

	return cmd
}
