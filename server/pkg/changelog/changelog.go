package changelog

import (
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/models"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func List(changelogDirPath string) ([]*models.Changelog, error) {

	out := []*models.Changelog{}

	dirEntries, err := ioutil.ReadDir(changelogDirPath)
	if err != nil {
		return nil, err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		date, err := time.Parse("2006-01-02", strings.TrimSuffix(dirEntry.Name(), ".md"))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse filename to YYYY-MM-DD date %s", dirEntry.Name())
		}

		dat, err := os.ReadFile(path.Join(changelogDirPath, dirEntry.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s", dirEntry.Name())
		}

		out = append(out, &models.Changelog{Date: date, Content: string(dat)})
	}
	return out, nil
}
