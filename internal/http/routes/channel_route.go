package routes

import (
	"github.com/danyouknowme/smthng/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

type channelRoutes struct {
	router         *gin.RouterGroup
	channelHandler handlers.ChannelHandler
	authMiddleware gin.HandlerFunc
}

func NewChannelRoutes(router *gin.RouterGroup, channelHandler handlers.ChannelHandler, authMiddleware gin.HandlerFunc) *channelRoutes {
	return &channelRoutes{
		router:         router,
		channelHandler: channelHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *channelRoutes) Register() {
	channelRoutes := r.router.Group("/channels")
	channelRoutes.Use(r.authMiddleware)
	{
		channelRoutes.POST("", r.channelHandler.CreateNewChannel)
	}
}
