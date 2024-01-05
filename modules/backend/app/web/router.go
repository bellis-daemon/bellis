package web

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-gonic/gin"
)

// ServeWeb serves the gRPC and HTTP endpoints using the provided net.Listener.
// It wraps the gRPC server, sets up routing for callback services, and starts serving requests using the gin router.
func ServeWeb(ctx context.Context, lis net.Listener) {
	router := gin.Default()
	apiRouter := router.Group("api")
	{
		apiRouter.GET("ip", services.GetIpInfo())
		callbackRouter := apiRouter.Group("callback")
		{
			callbackRouter.POST("telegram", services.TelegramCallbackService())
		}
		chartsRouter := apiRouter.Group("charts")
		{
			chartsRouter.GET(":id/request-time.png", services.RequestTimeChart())
		}
		sentrySingletonRouter:=apiRouter.Group("sentry-singleton")
		{
			sentrySingletonRouter.POST("refresh",services.SentrySingletonRefresh())
		}
	}

	srv := &http.Server{
		Handler: router,
	}
	go func() {
		if err := srv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	<-ctx.Done()
	srv.Shutdown(context.Background())
}
