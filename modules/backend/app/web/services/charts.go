package services

import (
	"fmt"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/gin-gonic/gin"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"github.com/wcharczuk/go-chart/v2"
	"net/http"
	"time"
)

func RequestTimeChart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		glgf.Debug(id)
		query, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(`
from(bucket: "backend")
  |> range(start: -15m)
  |> filter(fn: (r) => r["_field"] == "c_response_time")
  |> filter(fn: (r) => r["id"] == "%s")`, id))
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		var values []float64
		var times []time.Time
		for query.Next() {
			values = append(values, cast.ToFloat64(query.Record().Value()))
			times = append(times, query.Record().Time())
		}
		glgf.Debug(values)
		series := chart.TimeSeries{
			Name: "Request times",
			Style: chart.Style{
				StrokeColor: chart.ColorBlue,
				FillColor:   chart.ColorBlue.WithAlpha(100),
			},
			XValues: times,
			YValues: values,
		}
		graph := chart.Chart{
			Width:  1280,
			Height: 720,
			Background: chart.Style{
				Padding: chart.Box{
					Top: 50,
				},
			},
			YAxis: chart.YAxis{
				Name: "Elapsed Millis",
				TickStyle: chart.Style{
					TextRotationDegrees: 45.0,
				},
				ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%d ms", int(v.(float64)))
				},
			},
			XAxis: chart.XAxis{
				ValueFormatter: chart.TimeMinuteValueFormatter,
				GridMajorStyle: chart.Style{
					StrokeColor: chart.ColorAlternateGray,
					StrokeWidth: 1.0,
				},
			},
			Series: []chart.Series{
				series,
				chart.LastValueAnnotationSeries(series),
			},
		}
		graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}

		err = graph.Render(chart.PNG, ctx.Writer)
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Writer.Flush()
		ctx.Status(http.StatusOK)
	}
}
