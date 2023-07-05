package routes

import (
	"github.com/danyouknowme/smthng/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, userHandler handlers.UserHandler) {
	userRoutes := router.Group("/auth")
	{
		userRoutes.POST("/register", userHandler.Register)
	}
}
