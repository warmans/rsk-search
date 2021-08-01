package data

import (
	"encoding/json"
	"fmt"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
)

func ImportTranscriptsCmd() *cobra.Command {

	var inputDir string

	cmd := &cobra.Command{
		Use:   "transcribe",
		Short: "Import transcribed data from pilkipedia",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			for dateStr, name := range meta.XfmEpisodeNames() {

				lg := logger.With(zap.String("name", name), zap.String("publication", meta.PublicationXFM), zap.String("date", dateStr))

				date, err := time.Parse(time.RFC3339, dateStr)
				if err != nil {
					return err
				}
				transcript, err := loadTranscript(inputDir, meta.PublicationXFM, date)
				if err != nil {
					return nil
				}
				if transcript == nil {
					lg.Info("No transcript")
					continue
				}

				ep, err := data.LoadEpisode(cfg.dataDir, meta.PublicationXFM, name)
				if err != nil {
					return err
				}
				if ep == nil {
					lg.Info("No episode was initialized for this transcript")
					continue
				}

				// include the whole transcript
				ep.Transcript = transcript.Transcript

				// add any metadata provided by the transcripts
				if ep.Meta == nil {
					ep.Meta = models.Metadata{}
				}
				for k, v := range transcript.Meta {
					ep.Meta[k] = v
				}

				if err := data.ReplaceEpisodeFile(cfg.dataDir, ep); err != nil {
					return err
				}
				lg.Info("ok!")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./script/pilkipedia-scraper/raw", "Path to raw scraped files")

	return cmd
}

func loadTranscript(inputDir string, publication string, date time.Time) (*models.Transcript, error) {

	f, err := os.Open(path.Join(inputDir, fmt.Sprintf("transcript-%s-%s.json", publication, util.ShortDate(date))))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	e := &models.Transcript{}

	dec := json.NewDecoder(f)
	return e, dec.Decode(e)
}
