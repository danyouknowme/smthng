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
	"github.com/danyouknowme/smthng/internal/config"
	"github.com/danyouknowme/smthng/internal/http/routes"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer *http.Server
	config     *config.AppConfig
}

func NewApp(config *config.AppConfig) *App {
	router := initRouter()

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	return &App{
		httpServer: server,
		config:     config,
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

func initRouter() *gin.Engine {
	router := gin.New()

	hub := ws.NewWebsocketHub(&ws.Config{
		Redis: nil,
	})
	go hub.Run()

	routes.SetupWebSocketRoutes(router, hub)

	return router
}
