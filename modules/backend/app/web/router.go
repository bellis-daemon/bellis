package web

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// ServeWeb serves the gRPC and HTTP endpoints using the provided net.Listener.
// It wraps the gRPC server, sets up routing for callback services, and starts serving requests using the gin router.
func ServeWeb(ctx context.Context, lis net.Listener) {
	exec.Command(
		"/headless-shell/headless-shell",
		"--no-sandbox",
		"--use-gl=angle",
		"--use-angle=swiftshader",
		"--remote-debugging-address=0.0.0.0",
		"--remote-debugging-port=9222").
		Start()

	store := persistence.NewInMemoryStore(time.Minute)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	apiRouter := router.Group("api")
	{
		apiRouter.GET("ip", services.GetIpInfo())
		callbackRouter := apiRouter.Group("callback")
		{
			callbackRouter.POST("telegram", services.TelegramCallbackService())
		}
		chartsRouter := apiRouter.Group("charts")
		{
			chartsRouter.GET(":id/response-time.html", cache.CachePage(store, time.Minute, services.ResponseTimeChart(services.ResponseTimeChartModeHtml)))
			chartsRouter.GET(":id/response-time.png", cache.CachePage(store, time.Minute, services.ResponseTimeChart(services.ResponseTimeChartModePng)))
		}
		sentrySingletonRouter := apiRouter.Group("sentry-singleton")
		{
			sentrySingletonRouter.POST("refresh", services.SentrySingletonRefresh())
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
