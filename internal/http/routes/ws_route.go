package routes

import (
	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/danyouknowme/smthng/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type wsRoutes struct {
	router         *gin.RouterGroup
	hub            *ws.Hub
	jwtService     jwt.JWTService
	authMiddleware gin.HandlerFunc
}

func NewWebSocketRoutes(router *gin.RouterGroup, hub *ws.Hub, jwtService jwt.JWTService, authMiddleware gin.HandlerFunc) *wsRoutes {
	return &wsRoutes{
		router:         router,
		hub:            hub,
		jwtService:     jwtService,
		authMiddleware: authMiddleware,
	}
}

func (r *wsRoutes) Register() {
	wsRoutes := r.router.Group("/ws")
	wsRoutes.Use(r.authMiddleware)
	{
		wsRoutes.GET("", func(c *gin.Context) {
			ws.ServeWs(r.hub, c)
		})
	}
}
