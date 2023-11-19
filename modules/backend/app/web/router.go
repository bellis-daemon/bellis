package web

import (
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/minoic/glgf"
	"net"
)

func ServeWeb(lis net.Listener) {
	wrappedGrpc := grpcweb.WrapServer(mobile.Server())
	router := gin.Default()
	router.Use(func(context *gin.Context) {
		if wrappedGrpc.IsGrpcWebRequest(context.Request) || context.Request.Method == "OPTIONS" {
			glgf.Debug("gprc request: ", context.Request.URL)
			gin.WrapH(wrappedGrpc)
			context.Abort()
			return
		}
		glgf.Debug("http request: ", context.Request.URL)
	})
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
