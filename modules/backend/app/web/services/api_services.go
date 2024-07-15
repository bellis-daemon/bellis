package services

import (
	"net/http"
	"time"

	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/geo"
	"github.com/gin-gonic/gin"
	"github.com/minoic/glgf"
)

func GetIpInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		res, err := geo.FromLocal(ip)
		if err != nil {
			glgf.Error(err)
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusOK, res.String())
	}
}

func GetPingInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ret := map[string]any{
			"BuildTime": common.BuildTime,
			"GoVersion": common.GoVersion,
			"UnixMilli": time.Now().UnixMilli(),
			"Hostname":  common.Hostname(),
		}
		selfGeo, err := geo.Self()
		if err != nil {
			glgf.Error(err)
			ret["Region"] = "Unknown"
			ret["Country"] = "Unknown"
			ret["City"] = "Unknown"
			ret["ISP"] = "Unknown"
		} else {
			ret["Region"] = selfGeo.Region
			ret["Country"] = selfGeo.Country
			ret["City"] = selfGeo.City
			ret["ISP"] = selfGeo.ISP
		}
		ctx.JSON(http.StatusOK, ret)
	}
}
