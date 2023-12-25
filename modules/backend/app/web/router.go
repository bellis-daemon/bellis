package web

import (
	"net"

	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-gonic/gin"
)

// ServeWeb serves the gRPC and HTTP endpoints using the provided net.Listener.
// It wraps the gRPC server, sets up routing for callback services, and starts serving requests using the gin router.
func ServeWeb(lis net.Listener) {
	router := gin.Default()
	{
		callbackRouter := router.Group("callback")
		{
			callbackRouter.POST("telegram", services.TelegramCallbackService())
		}
		apiRouter := router.Group("api")
		{
			apiRouter.GET("ip", services.GetIpInfo())
		}
	}
	err := router.RunListener(lis)
	if err != nil {
		panic(err)
	}
}
