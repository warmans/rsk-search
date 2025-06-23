package chart

import (
	"context"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"slices"
	"strings"
)

type Kind string

const (
	RatingAvg       Kind = "avg"
	RatingCounts    Kind = "count"
	RatingBreakdown Kind = "breakdown"
)

func GenerateBreakdownChart(
	ctx context.Context,
	client api.TranscriptServiceClient,
	filterOrNil *filter.Filter,
) (*gg.Context, error) {
	defaultFilter := filter.Or(
		filter.Eq("publication_type", filter.String("radio")),
		filter.Eq("publication_type", filter.String("podcast")),
	)

	f := defaultFilter
	if filterOrNil != nil {
		f = filter.And(defaultFilter, *filterOrNil)
	}

	transcripts, err := client.ListTranscripts(ctx, &api.ListTranscriptsRequest{IncludeRatingBreakdown: true, Filter: filter.MustPrint(f)})
	if err != nil {
		return nil, err
	}

	allSeries := getAllAuthorSeries(transcripts)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 10})

	canvas := gg.NewContext(2000, 600)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	xScale := gochart.NewXScaleFromLabels(allSeries.episodeLabels)
	yScale := gochart.NewFixedYScale(10, 5.0)

	layout := gochart.NewDynamicLayout(
		gochart.NewStdYAxis(yScale),
		gochart.NewCompactXAxis(
			allSeries.episodeLabels,
			xScale,
			gochart.XCompactFontStyles(style.FontFace(face)),
		),
		append([]gochart.Plot{
			gochart.NewYGrid(yScale)},
			createPointPlots(yScale, xScale, allSeries.values)...,
		)...,
	)

	if err := layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas)); err != nil {
		return nil, err
	}

	return canvas, nil

}

func GenerateRatingsChart(
	ctx context.Context,
	client api.TranscriptServiceClient,
	filterOrNil *filter.Filter,
	author *string,
	sort bool,
	kind Kind,
) (*gg.Context, error) {

	defaultFilter := filter.Or(
		filter.Eq("publication_type", filter.String("radio")),
		filter.Eq("publication_type", filter.String("podcast")),
	)

	f := defaultFilter
	if filterOrNil != nil {
		f = filter.And(defaultFilter, *filterOrNil)
	}

	transcripts, err := client.ListTranscripts(ctx, &api.ListTranscriptsRequest{IncludeRatingBreakdown: true, Filter: filter.MustPrint(f)})
	if err != nil {
		return nil, err
	}

	var series gochart.Series
	var yScale gochart.YScale = gochart.NewFixedYScale(10, 5)

	if author != nil {
		series = createAuthorSeries(transcripts, *author)
	} else {
		switch kind {
		case RatingAvg:
			series = createAveragesSeries(transcripts)
		case RatingCounts:
			series = createCountSeries(transcripts)
			yScale = gochart.NewYScale(10, series)
		}

	}

	if sort {
		series = sortSeriesHighLow(series)
	}

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 10})

	canvas := gg.NewContext(2000, 600)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	xScale := gochart.NewXScale(series, 0)

	layout := gochart.NewDynamicLayout(
		gochart.NewStdYAxis(yScale),
		gochart.NewCompactXAxis(series.Xs(), xScale, gochart.XCompactFontStyles(style.FontFace(face))),
		append([]gochart.Plot{
			gochart.NewYGrid(yScale)},
			createAvgPlot(yScale, xScale, series),
			createBarPlot(yScale, xScale, series),
		)...,
	)

	if err := layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas)); err != nil {
		return nil, err
	}

	return canvas, nil
}

func createAveragesSeries(transcripts *api.TranscriptList) gochart.Series {
	XYs := struct {
		X []string
		Y []float64
	}{}

	for _, v := range transcripts.Episodes {
		XYs.X = append(XYs.X, v.ShortId)
		XYs.Y = append(XYs.Y, float64(v.RatingScore))
	}

	return gochart.NewXYSeries(XYs.X, XYs.Y)
}

func createCountSeries(transcripts *api.TranscriptList) gochart.Series {
	XYs := struct {
		X []string
		Y []float64
	}{}

	for _, v := range transcripts.Episodes {
		XYs.X = append(XYs.X, v.ShortId)
		XYs.Y = append(XYs.Y, float64(len(v.RatingBreakdown)))
	}

	return gochart.NewXYSeries(XYs.X, XYs.Y)
}

func createAuthorSeries(transcripts *api.TranscriptList, author string) gochart.Series {
	XYs := struct {
		X []string
		Y []float64
	}{}

	for _, v := range transcripts.Episodes {
		authorRating := 0.0
		if rating, ok := v.RatingBreakdown[fmt.Sprintf("discord:%s", author)]; ok {
			authorRating = float64(rating)
		}
		XYs.X = append(XYs.X, v.ShortId)
		XYs.Y = append(XYs.Y, authorRating)
	}

	return gochart.NewXYSeries(XYs.X, XYs.Y)
}

func createBarPlot(yScale gochart.YScale, xScale gochart.XScale, series gochart.Series) gochart.Plot {

	bar := gochart.NewBarsPlot(yScale, xScale, series, gochart.PlotStyleFn(func(v float64) style.Opts {
		if v <= 1 {
			return style.Opts{style.Color(color.RGBA{234, 85, 67, 255})}
		}
		if v <= 2 {
			return style.Opts{style.Color(color.RGBA{239, 156, 31, 255})}
		}
		if v <= 3 {
			return style.Opts{style.Color(color.RGBA{237, 224, 90, 255})}
		}
		if v <= 4 {
			return style.Opts{style.Color(color.RGBA{188, 207, 49, 255})}
		}
		return style.Opts{style.Color(color.RGBA{133, 187, 68, 255})}

	}))
	return bar
}

func createPointPlots(yScale gochart.YScale, xScale gochart.XScale, series [][]float64) []gochart.Plot {

	// the points need to be sized compared to others in the same X position
	// so count up how many are in each bucket per tick then use this in the
	// PointSizeFn
	weights := map[int]map[string]float64{}
	for _, label := range xScale.Labels() {
		weights[label.Tick] = make(map[string]float64)
		for _, ser := range series {
			if ser[label.Tick] == 0 {
				continue
			}
			weights[label.Tick][fmt.Sprintf("%0.1f", ser[label.Tick])]++
		}
	}

	plots := make([]gochart.Plot, len(series))
	for k, ser := range series {
		points := gochart.NewPointsPlot(
			yScale,
			xScale,
			gochart.NewYSeries(ser),
			gochart.PlotStyleFn(func(v float64) style.Opts {
				if v <= 1 {
					return style.Opts{style.Color(color.RGBA{234, 85, 67, 255})}
				}
				if v <= 2 {
					return style.Opts{style.Color(color.RGBA{239, 156, 31, 255})}
				}
				if v <= 3 {
					return style.Opts{style.Color(color.RGBA{237, 224, 90, 255})}
				}
				if v <= 4 {
					return style.Opts{style.Color(color.RGBA{188, 207, 49, 255})}
				}
				return style.Opts{style.Color(color.RGBA{133, 187, 68, 255})}

			}),
			gochart.PointSizeFn(func(v float64, label gochart.Label) float64 {
				maxWeight := 0.0
				for _, weight := range weights[label.Tick] {
					if weight > maxWeight {
						maxWeight = weight
					}
				}
				// if there is only one vote then just make the dot small
				if weights[label.Tick][fmt.Sprintf("%0.1f", v)] == 1 {
					return 2
				}
				// otherwise scale it to a proportional size
				return scaleBetween(weights[label.Tick][fmt.Sprintf("%0.1f", v)], 1, 7, 0, maxWeight)
			}))
		plots[k] = points
	}

	return plots
}

func scaleBetween(unscaledNum, minAllowed, maxAllowed, min, max float64) float64 {
	if unscaledNum == 0 {
		return 0
	}
	return (maxAllowed-minAllowed)*(unscaledNum-min)/(max-min) + minAllowed
}

func createAvgPlot(yScale gochart.YScale, xScale gochart.XScale, series gochart.Series) gochart.Plot {

	sum := 0.0
	for _, v := range series.Ys() {
		sum += v
	}

	avg := sum / float64(len(series.Ys()))

	avgs := []float64{}
	for range series.Xs() {
		avgs = append(avgs, avg)
	}

	return gochart.NewLinesPlot(yScale, xScale, gochart.NewXYSeries(series.Xs(), avgs), gochart.PlotStyle(
		style.Color(color.RGBA{A: 255, R: 255}),
	))
}

func sortSeriesHighLow(series gochart.Series) gochart.Series {
	tuples := []struct {
		x string
		y float64
	}{}

	for k, v := range series.Xs() {
		tuples = append(tuples, struct {
			x string
			y float64
		}{x: v, y: series.Y(k)})
	}

	slices.SortFunc(
		tuples,
		func(a, b struct {
			x string
			y float64
		}) int {
			if a.y > b.y {
				return -1
			}
			if a.y < b.y {
				return 1
			}
			return 0
		},
	)

	newX := []string{}
	newY := []float64{}
	for _, v := range tuples {
		newX = append(newX, v.x)
		newY = append(newY, v.y)
	}

	return gochart.NewXYSeries(newX, newY)
}

type authorSeries struct {
	episodeLabels []string
	authorLabels  []string
	values        [][]float64
}

func getAllAuthorSeries(list *api.TranscriptList) authorSeries {

	authors := getUniqueAuthors(list)
	s := authorSeries{
		authorLabels:  authors,
		values:        make([][]float64, len(authors)),
		episodeLabels: make([]string, 0),
	}

	for _, episode := range list.Episodes {
		s.episodeLabels = append(s.episodeLabels, episode.ShortId)
		for authorKey, authorName := range authors {
			rating, ok := episode.RatingBreakdown[authorName]
			if ok {
				s.values[authorKey] = append(s.values[authorKey], float64(rating))
			} else {
				s.values[authorKey] = append(s.values[authorKey], 0)
			}
		}
	}

	return s
}

func getUniqueAuthors(list *api.TranscriptList) []string {
	unique := map[string]struct{}{}
	for _, author := range list.Episodes {
		for authorName := range author.RatingBreakdown {
			if strings.HasPrefix(authorName, "discord:") {
				unique[authorName] = struct{}{}
			}
		}
	}

	out := []string{}
	for k := range unique {
		out = append(out, k)
	}
	return out
}
