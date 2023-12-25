package services

import (
	"net/http"

	"github.com/bellis-daemon/bellis/common/geo"
	"github.com/gin-gonic/gin"
	"github.com/minoic/glgf"
)

func GetIpInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		glgf.Debug(ip)
		res, err := geo.FromLocal(ip)
		if err != nil {
			glgf.Error(err)
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusOK, res.String())
	}
}
