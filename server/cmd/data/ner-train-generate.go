package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jdkato/prose/v2"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type RawLabeledEntity struct {
	Start int
	End   int
	Label string
	Text  string
}

type RawEntityContext struct {
	Accept bool
	Spans  []RawLabeledEntity
	Text   string
}

func NERTrainGenerateCmd() *cobra.Command {

	var inputDir string
	var outputFile string
	var numEpisodes int
	var skipEpisodes int

	cmd := &cobra.Command{
		Use:   "ner-train-generate",
		Short: "Extract some raw training data that must be manually updated to train the model",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: "+err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			reader := bufio.NewReader(os.Stdin)

			output, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			defer output.Close()

			dirEntries, err := ioutil.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for k, dirEntry := range dirEntries {
				if k < skipEpisodes {
					continue
				}

				logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

				episode := &models.Transcript{}
				if err := util.WithReadJSONFileDecoder(path.Join(inputDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				for _, v := range episode.Transcript {
					doc, err := prose.NewDocument(v.Content, prose.WithSegmentation(false))
					if err != nil {
						logger.Error("failed to parse text", zap.Error(err))
						continue
					}
					for _, ent := range doc.Entities() {
						entCtx := RawEntityContext{
							Text:   v.Content,
							Spans:  []RawLabeledEntity{},
							Accept: true,
						}
						pos := strings.Index(v.Content, ent.Text) // hopefully entities aren't duplicated significantly
						if pos == -1 {
							continue
						}
						entCtx.Spans = append(
							entCtx.Spans,
							RawLabeledEntity{Start: pos, End: pos + len(ent.Text), Label: ent.Label, Text: ent.Text},
						)

						fmt.Printf("%s: %s accept? ", ent.Label, ent.Text)
						r, _, err := reader.ReadLine()
						if err != nil {
							continue
						}
						entCtx.Accept = strings.TrimSpace(string(r)) == "y"

						bs, err := json.Marshal(entCtx)
						if err != nil {
							return err
						}

						bs = append(bs, []byte("\n")...)
						if _, err := output.Write(bs); err != nil {
							return err
						}
					}
				}

				numEpisodes--
				if numEpisodes == 0 {
					break
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw scraped files")
	cmd.Flags().StringVarP(&outputFile, "output-file", "o", "./var/data/ner/training-data-new.json", "Output file")
	cmd.Flags().IntVarP(&numEpisodes, "num-episodes", "n", 1, "Number of episodes to parse")
	cmd.Flags().IntVarP(&skipEpisodes, "skip-episodes", "s", 1, "Number of episodes to skip")

	return cmd
}
