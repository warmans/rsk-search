package data

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"go.uber.org/zap"
	"html/template"
	"os"
	"path"
	"strings"
)

var tmpl = template.Must(template.New("tmpl").Funcs(template.FuncMap{
	"actorColor": func(actor string, darken bool) string {
		switch strings.ToLower(actor) {
		case "ricky":
			if darken {
				return "#fffac9"
			}
			return "#fffdec"
		case "steve":
			if darken {
				return "#c9ffca"
			}
			return "#eeffec"
		case "karl":
			if darken {
				return "#c9fffc"
			}
			return "#ecffff"
		case "song":
			return "#ffecec"
		default:
			return "#ffffff"
		}
	},
	"typeColor": func(tp models.DialogType) string {
		switch tp {
		case models.DialogTypeSong:
			return "#ffecec"
		default:
			return "#ffffff"
		}
	},
}).Parse(`
<!DOCTYPE html>
<html>
  <head>
    <title>{{.ReleaseDate}}/Transcript</title>
  </head>
  <body>
	<p>This is a transcription of the {{.ReleaseDate}} episode, from {{.Publication}} Series {{.Series}}</p>
	{{ with .Transcript }}
		{{ range . }}
			{{  if eq .Type "chat" }}
			<div style="background-color: {{ actorColor .Actor false }};"><p>
				<span style="border: 0px solid {{ actorColor .Actor true }}; background-color:  {{ actorColor .Actor true }}; padding: 2px; font-weight: bold; margin-right: 5px;">{{ .Actor }}:</span> {{ .Content }}
			</p></div>
			{{ else }}
				<div style="background-color: {{ typeColor .Type }};"><p>
					<span>{{  if eq .Type "song" }}{{ .Type }}:{{ end }}</span> {{ .Content }}
				</p></div>
			{{ end }}
		{{ end }} 
	{{ end }}
    </table>
  </body>
</html>
`))

func GenerateHTMLCmd() *cobra.Command {

	var inputFile string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "generate-html",
		Short: "Convert a transcript to HTML (mediawiki compatible)",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			episode, err := data.LoadEpisodePath(inputFile)
			if err != nil {
				return err
			}

			outFile, err := os.Create(path.Join(outputPath, fmt.Sprintf("%s.html", episode.ID())))
			if err != nil {
				return err
			}
			defer func(outFile *os.File) {
				err := outFile.Close()
				if err != nil {
					logger.Error("failed to close outfile", zap.Error(err))
				}
			}(outFile)

			return tmpl.Execute(outFile, episode)
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "./var/data/episodes/ep-xfm-S2E08.json", "Path input JSON")
	cmd.Flags().StringVarP(&outputPath, "output-path", "o", "./var/data/export", "Path to output dir")

	return cmd
}
