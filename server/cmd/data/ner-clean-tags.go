package data

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"strings"
)

func NERCleanTagsCmd() *cobra.Command {

	var inputPath string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "ner-clean-tags",
		Short: "Removes any tags that are likely to be useless",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync()

			logger.Info("Importing transcript data from...", zap.String("path", inputPath))

			inputTags := meta.Tags{}

			err := util.WithReadJSONFileDecoder(inputPath, func(dec *json.Decoder) error {
				return dec.Decode(&inputTags)
			})
			if err != nil {
				return errors.Wrap(err, "failed to read tag data")
			}

			outputTags := meta.Tags{}

			for name, tag := range inputTags {
				if strings.ContainsAny(name, `…?")(-—.'`) {
					logger.Info(fmt.Sprintf("removed %s", name))
				} else {
					outputTags[name] = tag
				}
			}
			return util.WithCreateJSONFileEncoder(outputPath, func(enc *json.Encoder) error {
				return enc.Encode(outputTags)
			})
		},
	}

	cmd.Flags().StringVarP(&inputPath, "input-path", "i", "./pkg/meta/data/tags-new.json", "Path to input file")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "./pkg/meta/data/tags-clean.json", "Path to output file")

	return cmd
}
