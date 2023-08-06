package routes

import (
	"github.com/danyouknowme/smthng/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

type messageRoutes struct {
	router         *gin.RouterGroup
	messageHandler handlers.MessageHandler
	authMiddleware gin.HandlerFunc
}

func NewMessageRoutes(router *gin.RouterGroup, messageHandler handlers.MessageHandler, authMiddleware gin.HandlerFunc) *messageRoutes {
	return &messageRoutes{
		router:         router,
		messageHandler: messageHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *messageRoutes) Register() {
	messageRoutes := r.router.Group("/messages")
	messageRoutes.Use(r.authMiddleware)
	{
		messageRoutes.POST("/:channelID", r.messageHandler.CreateMessage)
		messageRoutes.PUT("/:messageID", r.messageHandler.EditMessage)
		messageRoutes.DELETE("/:messageID", r.messageHandler.DeleteMessage)
	}
}
