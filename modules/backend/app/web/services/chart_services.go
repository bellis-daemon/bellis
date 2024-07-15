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

		assetsHost := "https://cdn.jsdelivr.net/npm/echarts@5/dist/"

		switch mode {
		case ResponseTimeChartModeHtml:
		case ResponseTimeChartModeJpg:
			fallthrough
		case ResponseTimeChartModePng:
			req.Renderer = "canvas"
			req.Animation = false
			assetsHost = "./"
		case ResponseTimeChartModeSvg:
			req.Renderer = "svg"
			assetsHost = "./"
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
	|> aggregateWindow(every: 1m, fn: mean, createEmpty: true)
	|> fill(column: "_value", value: 0.0)`, id))
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		values := make([]opts.LineData, 1440)
		times := make([]string, 1440)
		var maxValue float64
		var lastValue float64
		for query.Next() {
			f := cast.ToFloat64(query.Record().Value())
			lastValue = f
			maxValue = math.Max(maxValue, f)
			val := cast.ToFloat64(query.Record().Value())
			tim := query.Record().Time().In(loc).Format("01/02 15:04")
			if len(times) != 0 && times[len(times)-1] == tim {
				continue
			}
			values = append(values, opts.LineData{Value: val})
			times = append(times, tim)
		}
		line := charts.NewLine()
		line.Animation = opts.Bool(req.Animation)
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle:  fmt.Sprintf("Bellis | Response Time - %s", entity.Name),
				Renderer:   req.Renderer,
				AssetsHost: assetsHost,
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
					Type: "line",
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
			charts.WithDataZoomOpts(opts.DataZoom{
				Type: "inside",
			}),
		)
		line.AddJSFuncs(`
const chart = %MY_ECHARTS%;
chart.setOption({
		graphic: 
		[
			{
				type: 'group',
				right: 0,
				bottom: 0,
				z: 100,
				onclick(){
				  	window.open("https://bellis.minoic.top","blank")
				},
				children: [{
						type: 'rect',
						left: 'center',
						top: 'center',
						z: 100,
						shape: {
							width: 300,
							height: 26
						},
						style: {
							fill: 'rgba(0,0,0,0.15)'
						}
					},
					{
						type: 'text',
						left: 'center',
						top: 'center',
						z: 100,
						style: {
							fill: '#fff',
							text: 'Chart By bellis.minoic.top 🌼',
							font: 'bold 14px sans-serif'
						}
					}
				]
			}, 
		],
	})`)
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
