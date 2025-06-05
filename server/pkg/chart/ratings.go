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
)

func GenerateRatingsChart(ctx context.Context, client api.TranscriptServiceClient, filterOrNil *filter.Filter, author *string) (*gg.Context, error) {

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

	var allSeries []gochart.Series
	if author != nil {
		allSeries = createAuthorSeries(transcripts, *author)
	} else {
		allSeries = createAveragesSeries(transcripts)
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

	yScale := gochart.NewYScale(10, allSeries[0])
	xScale := gochart.NewXScale(allSeries[0], 0)

	layout := gochart.NewDynamicLayout(
		gochart.NewStdYAxis(yScale),
		gochart.NewCompactXAxis(allSeries[0], xScale, gochart.XCompactFontStyles(style.FontFace(face))),
		append([]gochart.Plot{
			gochart.NewYGrid(yScale)},
			createBarPlots(yScale, xScale, allSeries)...,
		)...,
	)

	if err := layout.Render(canvas, gochart.BoundingBoxFromCanvas(canvas)); err != nil {
		return nil, err
	}

	return canvas, nil
}

func createAveragesSeries(transcripts *api.TranscriptList) []gochart.Series {
	XYs := struct {
		X []string
		Y []float64
	}{}

	for _, v := range transcripts.Episodes {
		XYs.X = append(XYs.X, v.ShortId)
		XYs.Y = append(XYs.Y, float64(v.RatingScore))
	}

	return []gochart.Series{gochart.NewXYSeries(XYs.X, XYs.Y)}
}

func createAuthorSeries(transcripts *api.TranscriptList, author string) []gochart.Series {
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

	return []gochart.Series{gochart.NewXYSeries(XYs.X, XYs.Y)}
}

func createBarPlots(yScale gochart.YScale, xScale gochart.XScale, series []gochart.Series) []gochart.Plot {
	plots := make([]gochart.Plot, len(series))
	for k, v := range series {
		bar := gochart.NewBarsPlot(yScale, xScale, v)
		bar.SetStyleFn(func(v float64) style.Opts {
			if v <= 1.5 {
				return style.Opts{style.Color(color.RGBA{R: 220, A: 255})}
			}
			if v > 1.5 && v < 2.5 {
				return style.Opts{style.Color(color.RGBA{R: 245, G: 138, B: 39, A: 255})}
			}
			return style.Opts{style.Color(color.RGBA{R: 23, G: 220, B: 0, A: 255})}

		})
		plots[k] = bar
	}
	return plots
}
