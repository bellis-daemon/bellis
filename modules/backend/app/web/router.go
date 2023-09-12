package web

import (
	"net"

	"github.com/gin-gonic/gin"
)

func ServeWeb(lis net.Listener) {
	router := gin.Default()
	router.RunListener(lis)
}
