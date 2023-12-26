package web

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

// ServeWeb serves the gRPC and HTTP endpoints using the provided net.Listener.
// It wraps the gRPC server, sets up routing for callback services, and starts serving requests using the gin router.
func ServeWeb(ctx context.Context, lis net.Listener) {
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
		chartsRouter := router.Group("charts")
		{
			chartsRouter.GET(":id/request-time.png", services.RequestTimeChart())
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
	select {
	case <-ctx.Done():
		srv.Shutdown(context.Background())
	}
}
