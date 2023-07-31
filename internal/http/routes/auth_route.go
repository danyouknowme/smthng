package routes

import (
	"github.com/danyouknowme/smthng/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	router      *gin.RouterGroup
	userHandler handlers.UserHandler
}

func NewAuthRoutes(router *gin.RouterGroup, userHandler handlers.UserHandler) *authRoutes {
	return &authRoutes{
		router:      router,
		userHandler: userHandler,
	}
}

func (r *authRoutes) Register() {
	authRoutes := r.router.Group("/auth")
	{
		authRoutes.POST("/register", r.userHandler.Register)
		authRoutes.POST("/login", r.userHandler.Login)
	}
}
