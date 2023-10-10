package web

import (
	"net"

	"github.com/gin-gonic/gin"
)

func ServeWeb(lis net.Listener) {
	router := gin.Default()
	err := router.RunListener(lis)
	if err != nil {
		panic(err)
	}
}
