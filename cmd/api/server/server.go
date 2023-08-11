package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/danyouknowme/smthng/internal/bussiness/usecases"
	"github.com/danyouknowme/smthng/internal/config"
	"github.com/danyouknowme/smthng/internal/datasources"
	"github.com/danyouknowme/smthng/internal/datasources/repositories"
	"github.com/danyouknowme/smthng/internal/http/handlers"
	"github.com/danyouknowme/smthng/internal/http/middleware"
	"github.com/danyouknowme/smthng/internal/http/routes"
	"github.com/danyouknowme/smthng/pkg/jwt"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer  *http.Server
	config      *config.AppConfig
	datasources datasources.DataSources
}

func NewApp(ds datasources.DataSources, cfg *config.AppConfig) *App {
	router := initRouter(ds, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return &App{
		httpServer:  server,
		config:      cfg,
		datasources: ds,
	}
}

func (a *App) Start() error {
	go func() {
		logger.Infof("Server listening at port %s", a.config.Port)
		logger.Info("Starting listening for HTTP requests completed")
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...")

	logger.Info("Unregistering datasources...")
	if err := a.datasources.Close(); err != nil {
		return fmt.Errorf("error when close datasources: %v", err)
	}
	logger.Info("Unregistering datasources completed")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown server: %v", err)
	}

	<-ctx.Done()
	logger.Info("Timeout of 3 seconds")
	logger.Info("Shutting down server completed")
	return nil
}

func initRouter(ds datasources.DataSources, cfg *config.AppConfig) *gin.Engine {
	router := gin.New()
	routeV1 := router.Group("/api/v1")

	jwtService := jwt.NewJWTService(cfg.TokenSymmetricKey)

	userRepository := repositories.NewUserRepository(ds)
	channelRepository := repositories.NewChannelRepository(ds)
	messageRepository := repositories.NewMessageRepository(ds)
	fileRepository := repositories.NewFileRepository(ds)
	sessionRepository := repositories.NewSessionRepository(ds)

	userUsecase := usecases.NewUserUsecase(userRepository)
	channelUsecase := usecases.NewChannelUsecase(channelRepository)
	messageUsecase := usecases.NewMessageUsecase(messageRepository, userRepository, fileRepository)
	sessionUsecase := usecases.NewSessionUsecase(sessionRepository)

	userHandler := handlers.NewUserHandler(userUsecase, sessionUsecase, jwtService)
	channelHandler := handlers.NewChannelHandler(channelUsecase)

	hub := ws.NewWebsocketHub(&ws.Config{
		ChannelUsecase: channelUsecase,
		Redis:          ds.GetRedisClient(),
	})
	go hub.Run()

	socketService := ws.NewSocketService(hub, channelRepository)
	messageHandler := handlers.NewMessageHandler(socketService, channelUsecase, messageUsecase)

	wsRoutes := routes.NewWebSocketRoutes(routeV1, hub, jwtService, middleware.AuthMiddleware(jwtService))
	wsRoutes.Register()

	channelRoutes := routes.NewChannelRoutes(routeV1, channelHandler, middleware.AuthMiddleware(jwtService))
	channelRoutes.Register()

	authRoutes := routes.NewAuthRoutes(routeV1, userHandler)
	authRoutes.Register()

	messageRoutes := routes.NewMessageRoutes(routeV1, messageHandler, middleware.AuthMiddleware(jwtService))
	messageRoutes.Register()

	return router
}
