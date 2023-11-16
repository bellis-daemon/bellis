package web

import (
	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/minoic/glgf"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeWeb(lis net.Listener) {
	router := gin.Default()
	{
		callbackRouter := router.Group("callback")
		{
			callbackRouter.POST("telegram", services.TelegramCallbackService())
		}
	}
	err := router.RunListener(lis)
	if err != nil {
		panic(err)
	}
}

func ServeGrpcWeb(lis net.Listener) {
	router := gin.Default()
	{
		router.Any("/*path", func(context *gin.Context) {
			glgf.Debug(context.Request.RemoteAddr, context.Request.RequestURI, context.Request.Proto, context.Request.Method)
			context.Status(http.StatusOK)
		})
	}
	err := router.RunListener(lis)
	if err != nil {
		panic(err)
	}
}
