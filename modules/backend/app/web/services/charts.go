package services

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/web/render"
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type responseTimeChartRequest struct {
	Animation bool   `form:"animation"`
	Timezone  string `form:"timezone"`
	Renderer  string `form:"renderer"`
}

type ResponseTimeChartMode uint

const (
	ResponseTimeChartModeHtml ResponseTimeChartMode = iota
	ResponseTimeChartModePng
	ResponseTimeChartModeJpg
	ResponseTimeChartModeSvg
)

func ResponseTimeChart(mode ResponseTimeChartMode) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := responseTimeChartRequest{
			Animation: true,
			Timezone:  "UTC",
			Renderer:  "canvas",
		}
		err := ctx.BindQuery(&req)
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		switch mode {
		case ResponseTimeChartModeHtml:
		case ResponseTimeChartModeJpg:
			fallthrough
		case ResponseTimeChartModePng:
			req.Renderer = "canvas"
			req.Animation = false
		case ResponseTimeChartModeSvg:
			req.Renderer = "svg"
		}

		loc, err := time.LoadLocation(req.Timezone)
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		id := ctx.Param("id")
		pid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		var entity models.Application
		err = storage.CEntity.FindOne(ctx, bson.M{"_id": pid}).Decode(&entity)
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		if !entity.Public.ExternalChart {
			ctx.AbortWithStatus(http.StatusPaymentRequired)
			return
		}
		query, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(`
from(bucket: "backend")
	|> range(start: -1d)
	|> filter(fn: (r) => r["_field"] == "c_response_time")
	|> filter(fn: (r) => r["id"] == "%s")
	|> aggregateWindow(every: 1m, fn: mean, createEmpty: false)`, id))
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		var values []opts.LineData
		var times []string
		var maxValue float64
		var lastValue float64
		for query.Next() {
			f := cast.ToFloat64(query.Record().Value())
			lastValue = f
			maxValue = math.Max(maxValue, f)
			values = append(values, opts.LineData{Value: cast.ToFloat64(query.Record().Value())})
			times = append(times, query.Record().Time().In(loc).Format("01/02 15:04:05"))
		}
		line := charts.NewLine()
		line.Animation = opts.Bool(req.Animation)
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle: fmt.Sprintf("Bellis | Response Time - %s", entity.Name),
				Renderer:  req.Renderer,
			}),
			charts.WithLegendOpts(opts.Legend{
				Type: "plain",
				Show: opts.Bool(true),
				Left: "right",
			}),
			charts.WithTitleOpts(opts.Title{
				Show:     opts.Bool(true),
				Title:    fmt.Sprintf("Response Time - %s (%.2f ms)", entity.Name, lastValue),
				Subtitle: fmt.Sprintf("Scheme: %s | CreatedAt: %s %s", entity.Scheme, entity.CreatedAt.In(loc).Format(time.DateTime), loc.String()),
			}),
			charts.WithTooltipOpts(opts.Tooltip{
				Show:           opts.Bool(true),
				Trigger:        "axis",
				TriggerOn:      "mousemove",
				ValueFormatter: string(opts.FuncOpts("(data)=>`${data.toFixed(2)} ms`")),
				AxisPointer: &opts.AxisPointer{
					Type:"line",
				},
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Max: math.Max(maxValue, 100.0),
				AxisLabel: &opts.AxisLabel{
					Show:         opts.Bool(true),
					ShowMaxLabel: opts.Bool(true),
					ShowMinLabel: opts.Bool(true),
					Formatter:    string(opts.FuncOpts("(data)=>`${data.toFixed(2)} ms`")),
				},
			}),
		)
		line.SetXAxis(times).AddSeries(entity.Name,
			values,
			charts.WithAreaStyleOpts(opts.AreaStyle{}),
		)
		switch mode {
		case ResponseTimeChartModeHtml:
			line.Render(ctx.Writer)
			ctx.Status(http.StatusOK)
			ctx.Writer.Header().Set("content-type", "text/html; charset=utf-8")
		case ResponseTimeChartModePng:
			line.BackgroundColor = "#FFFFFF"
			bts, err := render.MakeChartSnapshotPng(line.RenderContent())
			if err != nil {
				glgf.Error(err)
				ctx.Status(http.StatusInternalServerError)
				break
			}
			ctx.Writer.Write(bts)
			ctx.Status(http.StatusOK)
			ctx.Writer.Header().Set("content-type", "image/png")
		case ResponseTimeChartModeJpg:
			line.BackgroundColor = "#FFFFFF"
			bts, err := render.MakeChartSnapshotJpg(line.RenderContent())
			if err != nil {
				glgf.Error(err)
				ctx.Status(http.StatusInternalServerError)
				break
			}
			ctx.Writer.Write(bts)
			ctx.Status(http.StatusOK)
			ctx.Writer.Header().Set("content-type", "image/jpeg")
		case ResponseTimeChartModeSvg:
			bts, err := render.MakeChartSnapshotSvg(line.RenderContent())
			if err != nil {
				glgf.Error(err)
				ctx.Status(http.StatusInternalServerError)
				break
			}
			ctx.Writer.Write(bts)
			ctx.Status(http.StatusOK)
			ctx.Writer.Header().Set("content-type", "image/svg+xml")
		default:
			glgf.Error("invalid chart mode: ", mode)
			ctx.Status(http.StatusInternalServerError)
		}
	}
}
