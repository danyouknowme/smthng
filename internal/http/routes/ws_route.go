package routes

import (
	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/gin-gonic/gin"
)

func SetupWebSocketRoutes(router *gin.Engine, hub *ws.Hub) {
	wsRoutes := router.Group("/ws")
	wsRoutes.GET("", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})
}
