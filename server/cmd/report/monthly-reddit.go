package report

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/changelog"
	"github.com/warmans/rsk-search/pkg/store/common"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"go.uber.org/zap"
	"io"
	"os"
	"sort"
	"text/template"
	"time"
)

type contributionSummary struct {
	ContributionType string
	Points           float64
	NumContributions int64
}

type contributor struct {
	Contributor string
	Points      float64
}

type popularEpisode struct {
	Epid string
}

type importedAudio struct {
	Epid   string
	Epname string
}

func MonthlyRedditReport() *cobra.Command {

	dbCfg := &common.Config{}
	var changelogsDir string

	cmd := &cobra.Command{
		Use:   "reddit",
		Short: "reddit report",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
				}
			}()

			if dbCfg.DSN == "" {
				panic("dsn not set")
			}
			conn, err := rw.NewConn(dbCfg)
			if err != nil {
				return err
			}
			ctx := context.Background()

			var output = os.Stdout

			if _, err := fmt.Fprintf(output, "# Scrimpton Report %s %d", time.Now().Month().String(), time.Now().Year()); err != nil {
				return err
			}
			if err := changeLog(changelogsDir, output); err != nil {
				return err
			}

			err = conn.WithTx(func(tx *sqlx.Tx) error {
				if err := contibutions(ctx, tx, output); err != nil {
					return err
				}
				if err := topContributors(ctx, tx, output); err != nil {
					return err
				}
				if err := popularEpisodes(ctx, tx, output); err != nil {
					return err
				}
				if err := newTranscripts(ctx, tx, output); err != nil {
					return err
				}
				return nil
			})

			return err
		},
	}

	dbCfg.RegisterFlags(cmd.Flags(), "", "rw")
	cmd.Flags().StringVarP(&changelogsDir, "changelogs-dir", "i", "./var/changelogs", "Path to raw markdown files")

	return cmd
}

func changeLog(changelogsDirPath string, out io.Writer) error {

	ch, err := changelog.List(changelogsDirPath)
	if err != nil {
		return err
	}
	sort.Slice(ch, func(i, j int) bool {
		return ch[i].Date.After(ch[j].Date)
	})

	return template.Must(template.New("changelog").Parse(`

## Changelog 

{{ .Content }}`,
	)).Execute(out, ch[0])
}

func contibutions(ctx context.Context, conn *sqlx.Tx, out io.Writer) error {

	res, err := conn.QueryxContext(
		ctx,
		`SELECT contribution_type as type, SUM(points) as points, COUNT(1) as num_contribtions 
		FROM author_contribution ac 
		WHERE  date_trunc('month', ac.created_at) = $1 
		GROUP BY contribution_type`,
		fmt.Sprintf("%d-%d-01", time.Now().Year(), time.Now().Month()),
	)
	if err != nil {
		return err
	}
	defer res.Close()

	data := []contributionSummary{}
	for res.Next() {
		sum := contributionSummary{}
		if err := res.Scan(&sum.ContributionType, &sum.Points, &sum.NumContributions); err != nil {
			return err
		}
		data = append(data, sum)
	}

	return template.Must(template.New("contribution-count").Parse(`
## Contributions 

Type | Contributions | Points
:----|--------------:|------:
{{ range . }}{{ .ContributionType }} | {{ .NumContributions }} | {{ printf "%.2f" .Points }} 
{{ end }}`,
	)).Execute(out, data)
}

func topContributors(ctx context.Context, conn *sqlx.Tx, out io.Writer) error {

	res, err := conn.QueryxContext(
		ctx,
		`SELECT a.name, SUM(points) 
		FROM author_contribution ac LEFT JOIN author a ON ac.author_id = a.id 	
		WHERE  date_trunc('month', ac.created_at) = $1 
		GROUP BY a.name 
		ORDER BY SUM(points) DESC`,
		fmt.Sprintf("%d-%d-01", time.Now().Year(), time.Now().Month()),
	)
	if err != nil {
		return err
	}
	defer res.Close()

	data := []contributor{}
	for res.Next() {
		sum := contributor{}
		if err := res.Scan(&sum.Contributor, &sum.Points); err != nil {
			return err
		}
		data = append(data, sum)
	}

	return template.Must(template.New("contribution-count").Parse(`
## Top contributors

Author | Points
:------|:--------
{{ range . }}/u/{{ .Contributor }} | {{ printf "%.2f" .Points }} 
{{ end }}`,
	)).Execute(out, data)
}

func popularEpisodes(ctx context.Context, conn *sqlx.Tx, out io.Writer) error {

	res, err := conn.QueryxContext(
		ctx,
		`SELECT epid 
		FROM media_access_log mal 
		WHERE date_trunc('month', mal.time_bucket) = $1 AND media_type = 'episode' 
		GROUP BY epid 
		ORDER BY SUM(num_times_accessed) DESC
		LIMIT 10`,
		fmt.Sprintf("%d-%d-01", time.Now().Year(), time.Now().Month()),
	)
	if err != nil {
		return err
	}
	defer res.Close()

	data := []popularEpisode{}
	for res.Next() {
		sum := popularEpisode{}
		if err := res.Scan(&sum.Epid); err != nil {
			return err
		}
		data = append(data, sum)
	}

	return template.Must(template.New("contribution-count").Parse(`
## Most listened episodes

{{ range . }}1. [{{ .Epid }}](https://scrimpton.com/ep/{{ .Epid }}) 
{{ end }}`,
	)).Execute(out, data)
}

func newTranscripts(ctx context.Context, conn *sqlx.Tx, out io.Writer) error {

	res, err := conn.QueryxContext(
		ctx,
		`SELECT epid, epname FROM tscript_import ti WHERE date_trunc('month', created_at) = $1`,
		fmt.Sprintf("%d-%d-01", time.Now().Year(), time.Now().Month()),
	)
	if err != nil {
		return err
	}
	defer res.Close()

	data := []importedAudio{}
	for res.Next() {
		sum := importedAudio{}
		if err := res.Scan(&sum.Epid, &sum.Epname); err != nil {
			return err
		}
		data = append(data, sum)
	}

	return template.Must(template.New("contribution-count").Parse(`
## New files for transcription

{{ range . }} - {{ .Epid }}  
{{ end }}`,
	)).Execute(out, data)
}
