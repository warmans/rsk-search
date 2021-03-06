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

func TagDialogCmd() *cobra.Command {

	var inputDir string
	var modelPath string

	cmd := &cobra.Command{
		Use:   "tag-dialog",
		Short: "Dumps all the tags that are found in the transcripts",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			model := prose.ModelFromDisk(modelPath)

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

				// always rebuild episode tags
				episode.Tags = []string{}

				for k, v := range episode.Transcript {

					// update tags according to current tag metadata
					for oldText, oldTag := range episode.Transcript[k].ContentTags {
						tag := meta.Unalias(oldTag)
						if tag == "" || oldTag != tag {
							delete(episode.Transcript[k].ContentTags, oldText)
						}
						if tag != "" {
							episode.Transcript[k].ContentTags[oldText] = tag
						}
					}

					// add new tags
					doc, err := prose.NewDocument(v.Content, prose.WithSegmentation(false), prose.UsingModel(model))
					if err != nil {
						logger.Error("failed to parse text", zap.Error(err))
						continue
					}
					for _, ent := range doc.Entities() {
						canonicalTag := meta.Unalias(ent.Text)
						if canonicalTag == "" {
							continue
						}
						if episode.Transcript[k].ContentTags == nil {
							episode.Transcript[k].ContentTags = map[string]string{}
						}
						episode.Transcript[k].ContentTags[ent.Text] = canonicalTag
					}
					for _, t := range episode.Transcript[k].ContentTags {
						if !util.InStrings(t, episode.Tags...) {
							episode.Tags = append(episode.Tags, t)
						}
					}
				}
				if err := util.ReplaceEpisodeFile(inputDir, episode); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw scraped files")
	cmd.Flags().StringVarP(&modelPath, "model", "m", "./var/data/ner/rsk-model", "Model data")

	return cmd
}
