package data

import (
	"encoding/json"
	"fmt"
	"github.com/jdkato/prose/v2"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
)

func NERTestCmd() *cobra.Command {

	var inputFile string
	var modelPath string
	var defaultModel bool

	cmd := &cobra.Command{
		Use:   "ner-test",
		Short: "Output entities for a given file using the given model.",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: "+err.Error())
				}
			}()

			episode := &models.Transcript{}
			if err := util.WithReadJSONFileDecoder(inputFile, func(dec *json.Decoder) error {
				return dec.Decode(episode)
			}); err != nil {
				return err
			}

			opts := []prose.DocOpt{prose.WithSegmentation(false)}
			if !defaultModel {
				// Load our model, which we saved to a directory named "PRODUCT".
				model := prose.ModelFromDisk(modelPath)
				opts = append(opts, prose.UsingModel(model))
			}

			for _, v := range episode.Transcript {
				doc, err := prose.NewDocument(v.Content, opts...)
				if err != nil {
					return err
				}
				for _, ent := range doc.Entities() {
					fmt.Println(ent.Text, ent.Label)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "./var/data/episodes/ep-xfm-S2E08.json", "Path test file")
	cmd.Flags().StringVarP(&modelPath, "model", "m", "./var/data/ner/rsk-model", "Path to model")
	cmd.Flags().BoolVarP(&defaultModel, "default-model", "x", false, "use the default model instead")

	return cmd
}
