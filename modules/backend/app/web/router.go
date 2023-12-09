package web

import (
	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-gonic/gin"
	"net"
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
	}
	err := router.RunListener(lis)
	if err != nil {
		panic(err)
	}
}
