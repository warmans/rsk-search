package v2

import (
	"context"
	"fmt"
	"github.com/blugelabs/bluge"
	search2 "github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
	"github.com/blugelabs/bluge/search/highlight"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/bluge_query"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/store/ro"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
)

const (
	ResultContextLines = 3
	PageSize           = 10
)

func NewSearch(
	index *bluge.Reader,
	readOnlyDB *ro.Conn,
	episodeCache *data.EpisodeCache,
	audioUriPattern string,
	logger *zap.Logger,
) search.Searcher {
	return &Search{
		index:           index,
		readOnlyDB:      readOnlyDB,
		episodeCache:    episodeCache,
		logger:          logger,
		audioUriPattern: audioUriPattern,
	}
}

type Search struct {
	index           *bluge.Reader
	readOnlyDB      *ro.Conn
	episodeCache    *data.EpisodeCache
	logger          *zap.Logger
	audioUriPattern string
}

func (s *Search) Search(ctx context.Context, f filter.Filter, page int32) (*api.SearchResultList, error) {

	query, err := bluge_query.FilterToQuery(f)
	if err != nil {
		return nil, err
	}

	agg := aggregations.NewTermsAggregation(search2.Field("actor"), 25)
	agg.AddAggregation("transcript_id", aggregations.NewTermsAggregation(search2.Field("transcript_id"), 150))

	req := bluge.NewTopNSearch(PageSize, query).SetFrom(PageSize * int(page)).WithStandardAggregations()
	req.AddAggregation("actor_count_over_time", agg)

	dmi, err := s.index.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &api.SearchResultList{
		ResultCount: int32(dmi.Aggregations().Count()),
		Results:     []*api.SearchResult{},
		Stats:       map[string]*api.SearchStats{},
	}

	for _, actorBucket := range dmi.Aggregations().Aggregation("actor_count_over_time").(search2.BucketCalculator).Buckets() {
		res.Stats[actorBucket.Name()] = &api.SearchStats{
			Labels: []string{},
			Values: []float32{},
		}

		// fill in gaps in the episode stats to give a complete time-series
		for _, episodeID := range meta.EpisodeList() {
			var found = false
			for _, b := range actorBucket.Aggregation("transcript_id").(search2.BucketCalculator).Buckets() {
				if strings.HasPrefix(b.Name(), "ep-") && episodeID == strings.TrimPrefix(b.Name(), "ep-") {
					res.Stats[actorBucket.Name()].Labels = append(res.Stats[actorBucket.Name()].Labels, episodeID)
					res.Stats[actorBucket.Name()].Values = append(res.Stats[actorBucket.Name()].Values, float32(b.Count()))
					found = true
				}
			}
			if !found {
				res.Stats[actorBucket.Name()].Labels = append(res.Stats[actorBucket.Name()].Labels, episodeID)
				res.Stats[actorBucket.Name()].Values = append(res.Stats[actorBucket.Name()].Values, float32(0))
			}
		}
	}

	res.Results = []*api.SearchResult{}

	next, err := dmi.Next()
	for err == nil && next != nil {
		result := &api.SearchResult{
			Dialogs: []*api.DialogResult{},
		}
		err := s.readOnlyDB.WithStore(func(store *ro.Store) error {

			var innerErr error
			err = next.VisitStoredFields(func(field string, value []byte) bool {
				if innerErr != nil {
					return false
				}
				if field == "_id" {
					dialogs, episodeID, err := store.GetDialogWithContext(ctx, string(value), ResultContextLines)
					if err != nil {
						innerErr = err
						return false
					}
					// it's possible the episode won't be found, if so skip this result
					if episodeID == "" {
						s.logger.Warn("failed to find dialog line, is the DB out of sync with the index?", zap.String("dialog_id", string(value)))
						return true
					}
					ep, err := s.episodeCache.GetEpisode(episodeID)
					if err != nil {
						innerErr = errors.Wrapf(err, "episode ID: %s", episodeID)
						return false
					}
					result.Episode = ep.ShortProto(fmt.Sprintf(s.audioUriPattern, ep.ShortID()))

					// dialogs
					lines := make([]*api.Dialog, len(dialogs))
					for k, d := range dialogs {
						lines[k] = d.Proto(string(value) == d.ID)
					}
					result.Dialogs = append(result.Dialogs, &api.DialogResult{Transcript: lines, Score: float32(next.Score)})
				}
				return true
			})
			if err != nil {
				return fmt.Errorf("error accessing stored fields: %v", err)
			}
			if innerErr != nil {
				return fmt.Errorf("error processing result: %v", innerErr)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		if result.Episode != nil {
			res.Results = append(res.Results, result)
		}

		if next, err = dmi.Next(); err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Search) ListTerms(fieldName string, prefix string) (models.FieldValues, error) {

	terms := models.FieldValues{}

	fieldDict, err := s.index.DictionaryIterator(fieldName, nil, []byte(prefix), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := fieldDict.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	tfd, err := fieldDict.Next()
	for err == nil && tfd != nil && strings.TrimSpace(tfd.Term()) != "" {
		if !strings.HasPrefix(strings.ToLower(tfd.Term()), strings.ToLower(prefix)) {
			tfd, err = fieldDict.Next()
			continue
		}
		terms = append(terms, models.FieldValue{Value: tfd.Term(), Count: int32(tfd.Count())})
		if len(terms) > 500 {
			return nil, fmt.Errorf("too many terms for field '%s' returned (prefix: %s)", fieldName, prefix)
		}
		tfd, err = fieldDict.Next()
	}

	sort.Slice(terms, func(i, j int) bool {
		return terms[i].Count > terms[j].Count
	})

	return terms, nil
}

func (s *Search) PredictSearchTerms(ctx context.Context, prefix string, exact bool, numPredictions int32, f filter.Filter) (*api.SearchTermPredictions, error) {
	var q bluge.Query
	if f != nil {
		var prefixQuery filter.Filter
		if exact {
			prefixQuery = filter.Eq("content", filter.String(prefix))
		} else {
			prefixQuery = filter.Like("content", filter.String(prefix))
		}
		filterQuery, err := bluge_query.FilterToQuery(filter.And(f, prefixQuery))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create query")
		}
		q = filterQuery
	} else {
		if !exact {
			matchQuery := bluge.NewMatchQuery(prefix)
			matchQuery.SetField("content")
			q = matchQuery
		} else {
			matchQuery := bluge.NewMatchPhraseQuery(prefix)
			matchQuery.SetField("content")
			q = matchQuery
		}
	}

	// fetch extras in case some need to be discarded
	req := bluge.NewTopNSearch(maxInt(int(numPredictions), 50), q).SetFrom(0)
	req.IncludeLocations()

	highlighter := highlight.NewSimpleHighlighter(
		highlight.NewSimpleFragmenterSized(256),
		NewBracketFragmentFormatter(),
		"[...]",
	)

	dmi, err := s.index.Search(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "search failed")
	}

	predictions := &api.SearchTermPredictions{
		Prefix:      prefix,
		Predictions: []*api.Prediction{},
	}

	next, err := dmi.Next()
	for err == nil && next != nil {
		p := &api.Prediction{}
		err = next.VisitStoredFields(func(field string, value []byte) bool {
			switch field {
			case "content":
				p.Line = string(value)
				return false
			case "actor":
				p.Actor = string(value)
				return true
			}
			// id is in the format [epid]-[pos] e.g. ep-xfm-S1E06-347
			if field == "_id" {
				var err error
				p.Epid, p.Pos, err = extractEpidAndPos(string(value))
				if err != nil {
					s.logger.Error("error extracting data from index ID", zap.Error(err))
					return false
				}
			}
			return true
		})
		if err != nil {
			return nil, err
		}

		// de-duplicate results
		duplicate := false
		for _, v := range predictions.Predictions {
			if !duplicate {
				duplicate = !stringsAreNotTooSimilar(v.Line, p.Line)
			}
		}

		if !duplicate && stringsAreNotTooSimilar(prefix, p.Line) {
			// highlight and fragment result
			fragments := highlighter.BestFragments(next.Locations["content"], []byte(p.Line), 1)
			if len(fragments) > 0 {
				p.Fragment = fragments[0]
			}
			predictions.Predictions = append(predictions.Predictions, p)
		}

		if len(predictions.Predictions) >= int(numPredictions) {
			return predictions, nil
		}

		if next, err = dmi.Next(); err != nil {
			return nil, err
		}
	}

	return predictions, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func stringsAreNotTooSimilar(search, found string) bool {
	return strings.Trim(strings.ToLower(search), ".?,!") != strings.Trim(strings.ToLower(found), ".?,!")
}

// index id is in the format [epid]-[pos] e.g. ep-xfm-S1E06-347
func extractEpidAndPos(indexID string) (string, int32, error) {
	segments := strings.Split(indexID, "-")
	if len(segments) < 4 {
		return "", 0, fmt.Errorf("failed to extract data from index ID: %s", indexID)
	}
	posStr := segments[len(segments)-1]
	epid := strings.Join(segments[:len(segments)-1], "-")

	posInt, err := strconv.ParseInt(posStr, 10, 32)
	if err != nil {
		return "", 0, errors.Wrapf(err, "failed to parse pos int from ID: %s", indexID)
	}
	return epid, int32(posInt), nil
}
