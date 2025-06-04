package chart

import (
	"context"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/warmans/gochart"
	"github.com/warmans/gochart/pkg/style"
	"github.com/warmans/rsk-search/gen/api"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
)

func GenerateRatingsChart(ctx context.Context, client api.TranscriptServiceClient) (*gg.Context, error) {
	transcripts, err := client.ListTranscripts(ctx, &api.ListTranscriptsRequest{IncludeRatingBreakdown: true, Filter: `publication_type = "radio"`})
	if err != nil {
		return nil, err
	}

	allSeries := createAveragesSeries(transcripts)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 8})

	canvas := gg.NewContext(1200, 400)
	canvas.SetColor(color.White)
	canvas.DrawRectangle(0, 0, float64(canvas.Width()), float64(canvas.Height()))
	canvas.Fill()

	yScale := gochart.NewYScale(10, allSeries[0])
	xScale := gochart.NewXScale(allSeries[0], 0)

	layout := gochart.NewDynamicLayout(
		gochart.NewYAxis(yScale),
		gochart.NewXAxis(allSeries[0], xScale, gochart.XFontStyles(style.FontFace(face))),
		append([]gochart.Plot{
			gochart.NewYGrid(yScale)},
			createLinePlots(yScale, xScale, allSeries)...,
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

func createSeriesPerAuthor(transcripts *api.TranscriptList) []gochart.Series {
	uniqueRaterMap := map[string]struct{}{}
	for _, v := range transcripts.Episodes {
		for rater := range v.RatingBreakdown {
			uniqueRaterMap[rater] = struct{}{}
		}
	}

	// order will change each time it's rendered, but because the series are not really distinct, it doesn't matter.
	orderedRaters := make([]string, 0, len(uniqueRaterMap))
	for name := range uniqueRaterMap {
		orderedRaters = append(orderedRaters, name)
	}

	XYs := make([]struct {
		X []string
		Y []float64
	}, len(uniqueRaterMap))

	for _, v := range transcripts.Episodes {
		for raterIdx, name := range orderedRaters {
			rating, ok := v.RatingBreakdown[name]
			if !ok {
				rating = 0
			}
			XYs[raterIdx].X = append(XYs[raterIdx].X, v.ShortId)
			XYs[raterIdx].Y = append(XYs[raterIdx].Y, float64(rating))
		}
	}

	series := make([]gochart.Series, 0, len(orderedRaters))
	for _, v := range XYs {
		series = append(series, gochart.NewXYSeries(v.X, v.Y))
	}

	return series
}

func createPointPlots(yScale gochart.YScale, xScale gochart.XScale, series []gochart.Series) []gochart.Plot {
	plots := make([]gochart.Plot, len(series))
	for k, v := range series {
		plots[k] = gochart.NewPointsPlot(yScale, xScale, v)
	}
	return plots
}

func createLinePlots(yScale gochart.YScale, xScale gochart.XScale, series []gochart.Series) []gochart.Plot {
	plots := make([]gochart.Plot, len(series))
	for k, v := range series {
		plots[k] = gochart.NewLinesPlot(yScale, xScale, v)
	}
	return plots
}
