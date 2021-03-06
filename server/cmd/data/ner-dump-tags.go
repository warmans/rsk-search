package data

import (
	"encoding/json"
	"github.com/jdkato/prose/v2"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"path"
)

func NERDumpTagsCmd() *cobra.Command {

	var inputDir string
	var outputFile string
	var modelPath string

	cmd := &cobra.Command{
		Use:   "ner-dump-tags",
		Short: "Dumps all the tags that are found in the transcripts",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			model := prose.ModelFromDisk(modelPath)
			tags := make(meta.Tags)

			dirEntries, err := ioutil.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {

				logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

				episode := &models.Episode{}
				if err := util.WithReadJSONFileDecoder(path.Join(inputDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				for _, v := range episode.Transcript {
					doc, err := prose.NewDocument(v.Content, prose.WithSegmentation(false), prose.UsingModel(model))
					if err != nil {
						logger.Error("failed to parse text", zap.Error(err))
						continue
					}
					for _, ent := range doc.Entities() {
						if len(ent.Text) <= 1 {
							continue
						}
						if _, ok := tags[ent.Text]; ok {
							if !util.InStrings(ent.Label, tags[ent.Text].Kind...) {
								tags[ent.Text].Kind = append(tags[ent.Text].Kind, ent.Label)
							}
						} else {
							tags[ent.Text] = &meta.Tag{
								Kind: []string{ent.Label},
							}
						}
					}
				}
			}

			return util.WithCreateJSONFileEncoder(outputFile, func(enc *json.Encoder) error {
				return enc.Encode(tags)

			})
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw scraped files")
	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "./pkg/meta/data/tags-new.json", "Output file")
	cmd.Flags().StringVarP(&modelPath, "model", "m", "./var/data/ner/rsk-model", "Model data")

	return cmd
}
