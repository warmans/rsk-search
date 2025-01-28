package images

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/util"
	"log"
	"os"
	"path"
	"strings"
)

func ThumbsCmd() *cobra.Command {

	var inputDir string

	cmd := &cobra.Command{
		Use:   "thumbs",
		Short: "create thumbnails for images in directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := os.ReadDir(inputDir)
			if err != nil {
				return fmt.Errorf("failed to read dir: %w", err)
			}
			for _, v := range files {
				if v.IsDir() || !util.InStrings(strings.ToLower(path.Ext(v.Name())), ".jpg", ".jpeg", ".png", ".webp") {
					continue
				}
				if strings.HasSuffix(strings.TrimSuffix(v.Name(), path.Ext(v.Name())), ".thumb") {
					continue
				}
				src, err := imaging.Open(path.Join(inputDir, v.Name()))
				if err != nil {
					log.Fatalf("failed to open image: %v", err)
				}

				dst := imaging.Resize(src, 300, 0, imaging.Lanczos)

				err = imaging.Save(dst, util.ThumbPath(inputDir, v.Name()))
				if err != nil {
					return fmt.Errorf("failed to save image: %w", err)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/archive", "Path to images")

	return cmd
}
