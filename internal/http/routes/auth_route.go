package routes

import (
	"github.com/danyouknowme/smthng/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, userHandler handlers.UserHandler) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", userHandler.Register)
		authRoutes.POST("/login", userHandler.Login)
	}
}
