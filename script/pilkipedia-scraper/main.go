package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// when recording dialog position leave gaps to allow missing entries to be inserted
const positionGap = int64(100)

func main() {

	indexer := colly.NewCollector(
		colly.AllowedDomains("web.archive.org"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./cache/archive_org_cache"),
	)

	episodeDetailsCollector := indexer.Clone()

	indexer.OnHTML(`li > a`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains "coursera.org/learn"
		if strings.HasSuffix(e.Text, "/Transcript") {
			episodeDetailsCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})

	// per page scraper
	episodeDetailsCollector.OnHTML("div[id=content]", func(e *colly.HTMLElement) {

		episode := models.Episode{
			Transcript: []models.Dialog{},
			Meta:       []models.Metadata{},
		}

		fmt.Println("Loaded page ", e.Request.URL)
		episode.Source = e.Request.URL.String()

		var pageTitle *colly.HTMLElement
		e.ForEach("h1#firstHeading", func(i int, element *colly.HTMLElement) {
			pageTitle = element
		})

		// episode description should always be in the first p of the content.
		var pageDescription *colly.HTMLElement
		e.ForEach(".mw-parser-output > p:nth-child(1), #mw-content-text > p:nth-child(1)", func(i int, element *colly.HTMLElement) {
			pageDescription = element
		})

		if pageTitle != nil || pageDescription != nil {
			fmt.Println("Parsing meta...")
			meta, err := ParseMeta(pageTitle, pageDescription)
			if err != nil {
				fmt.Printf("Failed to parse meta: %s", err.Error())
				return
			}
			episode.Meta = meta
		}

		position := positionGap
		e.ForEach("#mw-content-text > div[style], .mw-parser-output > div[style]", func(i int, element *colly.HTMLElement) {
			dialog, err := ParseDialog(element)
			if err != nil {
				fmt.Printf("Failed to parse line: %s", err.Error())
				return
			}
			dialog.Position = position
			position += positionGap
			episode.Transcript = append(episode.Transcript, *dialog)
		})

		fName := fmt.Sprintf("./raw/transcript-%s.json", episode.CanonicalName())
		file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

		// if the file already exists add a numeric suffix up until 5 then fail as something is probably
		// broken.
		if err != nil {
			if err != os.ErrExist {
				log.Fatalf("Cannot create file %q: %s\n", fName, err)
			}
			suffix := 1
			for {
				file, err = os.OpenFile(strings.ReplaceAll(fName, ".json", fmt.Sprintf("-%s.json", suffix)), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
				if err != nil {
					if err == os.ErrExist {
						suffix++
						if suffix >= 5 {
							log.Fatalf("Refusing to create more than 5 files with duplicate names")
						}
					} else {
						log.Fatalf("failed to open file: %s", err)
					}
				} else {
					break
				}
			}
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(true)

		if err := enc.Encode(episode); err != nil {
			log.Fatalf("Failed to encode JSON: %s\n", err)
		}
	})

	if err := indexer.Visit("https://web.archive.org/web/20200704135748/http://www.pilkipedia.co.uk/wiki/index.php?title=Category:Transcripts"); err != nil {
		log.Fatalf("failed visit top level URL: %s", err)
	}
}

// e.g. This is a transcription of the 15 November 2003 episode, from Xfm Series 3
func ParseMeta(pageTitle *colly.HTMLElement, firstParagraph *colly.HTMLElement) ([]models.Metadata, error) {

	if pageTitle == nil && firstParagraph == nil {
		return nil, nil
	}

	meta := []models.Metadata{}

	date, publication := getRawMetaParts(firstParagraph)
	if date == "" && pageTitle != nil {
		// fall back to title
		date = strings.TrimSpace(strings.TrimSuffix(pageTitle.Text, "/Transcript"))
	}
	if date == "" && publication == "" {
		return nil, fmt.Errorf("couldn't parse meta from line: %s", firstParagraph.Text)
	}

	dateMeta := models.Metadata{
		Type:  models.MetadataTypeDate,
		Value: "",
	}

	// e.g.  15 November 2003
	parsed, err := time.Parse("02 January 2006", date)
	if err == nil {
		dateMeta.Value = parsed.Format(time.RFC3339)
	}

	meta = append(meta, dateMeta)

	// Xfm Series 3
	publication, series := parsePublication(publication)
	if publication != "" {
		meta = append(meta, models.Metadata{
			Type:  models.MetadataTypePublication,
			Value: publication,
		})
	}
	if series != "" {
		meta = append(meta, models.Metadata{
			Type:  models.MetadataTypeSeries,
			Value: series,
		})
	}

	return meta, nil
}

// should return [date, publication series N]
func getRawMetaParts(e *colly.HTMLElement) (string, string) {
	if e == nil {
		return "", ""
	}
	// try with tags
	texts := trimStrings(e.ChildTexts("a"))
	if len(texts) == 2 {
		return texts[0], texts[1]
	}
	// try with regex
	texts = trimStrings(regexp.MustCompile(`([0-9]{2}.+\w.+[0-9]{4}).+from(.+)`).FindStringSubmatch(e.Text))
	if len(texts) == 3 {
		return texts[1], texts[2]
	}
	return "", ""
}

func trimStrings(ss []string) []string {
	for k := range ss {
		ss[k] = strings.TrimSpace(ss[k])
	}
	return ss
}

func parsePublication(line string) (string, string) {
	parts := strings.Split(strings.ToLower(line), "series")
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func ParseDialog(el *colly.HTMLElement) (*models.Dialog, error) {

	content, contentPrefix := cleanContent(el)

	dialog := &models.Dialog{
		ID:      shortuuid.New(),
		Actor:   strings.ToLower(strings.TrimSuffix(strings.TrimSpace(el.ChildText("span")), ":")),
		Type:    models.DialogTypeUnkown,
		Content: content,
	}
	if contentPrefix == "song" {
		dialog.Type = models.DialogTypeSong
	} else {
		if dialog.Actor != "" {
			dialog.Type = models.DialogTypeChat
		}
	}

	return dialog, nil
}

func cleanContent(el *colly.HTMLElement) (string, string) {
	raw := strings.ReplaceAll(strings.TrimSpace(el.Text), "\n", "")
	parts := strings.SplitN(raw, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1]), strings.TrimSpace(strings.ToLower(parts[0]))
	}
	return raw, ""
}
