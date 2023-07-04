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
	"github.com/danyouknowme/smthng/internal/http/routes"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	router := initRouter()

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	return &App{
		httpServer: server,
	}
}

func (a *App) Start() error {
	go func() {
		logger.Info("Server listening at port 8000")
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
