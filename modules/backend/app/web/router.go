package web

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/bellis-daemon/bellis/common/openobserve"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	gin_cache "github.com/bellis-daemon/bellis/modules/backend/midwares/gin-cache"
	"github.com/bellis-daemon/bellis/modules/backend/midwares/gin-cache/persist"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// ServeWeb serves the gRPC and HTTP endpoints using the provided net.Listener.
// It wraps the gRPC server, sets up routing for callback services, and starts serving requests using the gin router.
func ServeWeb(ctx context.Context, lis net.Listener) {
	store := persist.NewRedisStore(storage.Redis())

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	openobserve.RegisterGin(router)
	apiRouter := router.Group("api")
	{
		apiRouter.GET("ip", services.GetIpInfo())
		apiRouter.GET("ping", services.GetPingInfo())
		callbackRouter := apiRouter.Group("callback")
		{
			callbackRouter.POST("telegram", services.TelegramCallbackService())
		}
		chartsRouter := apiRouter.Group("charts", gin_cache.CacheByRequestURI(store, time.Minute, gin_cache.WithPrefixKey("GIN_CACHE_")))
		{
			chartsRouter.GET(":id/response-time.html ", services.ResponseTimeChart(services.ResponseTimeChartModeHtml))
			chartsRouter.GET(":id/response-time.png", services.ResponseTimeChart(services.ResponseTimeChartModePng))
			chartsRouter.GET(":id/response-time.jpg", services.ResponseTimeChart(services.ResponseTimeChartModeJpg))
			chartsRouter.GET(":id/response-time.svg", services.ResponseTimeChart(services.ResponseTimeChartModeSvg))
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
