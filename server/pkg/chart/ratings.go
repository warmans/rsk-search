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
)

type Kind string

const (
	RatingAvg    Kind = "avg"
	RatingCounts Kind = "count"
)

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
		gochart.NewCompactXAxis(series, xScale, gochart.XCompactFontStyles(style.FontFace(face))),
		append([]gochart.Plot{
			gochart.NewYGrid(yScale)},
			createBarPlots(yScale, xScale, []gochart.Series{series})...,
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

func createBarPlots(yScale gochart.YScale, xScale gochart.XScale, series []gochart.Series) []gochart.Plot {
	plots := make([]gochart.Plot, len(series))
	for k, v := range series {
		bar := gochart.NewBarsPlot(yScale, xScale, v)
		bar.SetStyleFn(func(v float64) style.Opts {
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

		})
		plots[k] = bar
	}
	return plots
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
