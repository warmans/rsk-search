package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"log"
	"os"
	"regexp"
	"strconv"
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

		episode := models.Transcript{
			Episode:    -1,
			Series:     -1,
			Transcript: []models.Dialog{},
			Meta:       models.Metadata{},
		}

		fmt.Println("Loaded page ", e.Request.URL)

		episode.Meta[models.MetadataTypePilkipediaURL] = e.Request.URL.String()

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

			releaseDate, err := ParseReleaseDate(pageTitle, pageDescription)
			if err != nil {
				fmt.Printf("Failed to parse release date for %s: %s \n", pageTitle.Text, err.Error())
				releaseDate = &time.Time{}
			}
			episode.ReleaseDate = *releaseDate
			episode.Publication, episode.Series = ParsePublication(pageDescription)
		} else {
			fmt.Printf("No source for metadata found - skipping")
			return
		}

		position := positionGap
		e.ForEach("#mw-content-text > div[style], .mw-parser-output > div[style]", func(i int, element *colly.HTMLElement) {

			parseLine := func(el *colly.HTMLElement) bool {
				dialog, err := ParseDialog(el)
				if dialog == nil {
					fmt.Println("no dialog for: ", el.Text)
					return false
				}
				if err != nil {
					fmt.Printf("Failed to parse line: %s", err.Error())
					return false
				}
				dialog.Position = position
				position += positionGap
				episode.Transcript = append(episode.Transcript, *dialog)
				return true
			}

			parseLine(element)

			// at least one transcript has an unclosed tag so it needs to go a level deeper
			element.ForEach("div[style]", func(i int, unclosedDivTag *colly.HTMLElement) {
				parseLine(unclosedDivTag)
			})
		})

		if episode.Publication == "" || episode.ReleaseDate.IsZero() {
			fmt.Printf("Skipping episode with insufficient metadata - publication: %s date: %s\n", episode.Publication, episode.ReleaseDate.String())
			return
		}

		fName := fmt.Sprintf("./raw/transcript-%s-%s.json", episode.Publication, util.ShortDate(episode.ReleaseDate))
		file, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

		// if the file already exists add a numeric suffix up until 5 then fail as something is probably
		// broken.
		if err != nil {
			log.Fatal("failed to open file for writing ", err.Error())
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
func ParseReleaseDate(pageTitle *colly.HTMLElement, firstParagraph *colly.HTMLElement) (*time.Time, error) {

	if pageTitle == nil && firstParagraph == nil {
		return nil, nil
	}

	date, _ := getRawMetaParts(firstParagraph)
	if date == "" && pageTitle != nil {
		// fall back to title
		date = strings.TrimSpace(strings.TrimSuffix(pageTitle.Text, "/Transcript"))
	}
	if date == "" {
		return nil, fmt.Errorf("couldn't parse meta from line: %s", firstParagraph.Text)
	}

	// e.g.  15 November 2003
	parsed, err := time.Parse("02 January 2006", date)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func ParsePublication(firstParagraph *colly.HTMLElement) (string, int32) {
	_, publicationStr := getRawMetaParts(firstParagraph)
	return parsePublication(publicationStr)
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

func parsePublication(line string) (string, int32) {
	parts := strings.Split(strings.ToLower(line), "series")
	if len(parts) != 2 {
		return "", 0
	}

	publication := strings.TrimSpace(parts[0])

	series, err := strconv.Atoi(parts[1])
	if err != nil {
		return publication, -1
	}

	return publication, int32(series)
}

func ParseDialog(el *colly.HTMLElement) (*models.Dialog, error) {

	var dialog *models.Dialog

	el.ForEach("p", func(i int, pTag *colly.HTMLElement) {

		content, contentPrefix := cleanContent(pTag)

		dialog = &models.Dialog{
			ID:      shortuuid.New(),
			Actor:   strings.ToLower(strings.TrimSuffix(strings.TrimSpace(pTag.ChildText("span")), ":")),
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
	})
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
