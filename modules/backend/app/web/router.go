package web

import (
	"net"

	"github.com/bellis-daemon/bellis/modules/backend/app/web/services"
	"github.com/gin-gonic/gin"
)

func ServeWeb(lis net.Listener) {
	router := gin.Default()
	{
		callbackRouter:= router.Group("callback")
		{
			callbackRouter.POST("telegram",services.TelegramCallbackService())
		}
	}
	err := router.RunListener(lis)
	if err != nil {
		panic(err)
	}
}
