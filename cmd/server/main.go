package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danyouknowme/smthng/cmd/ws"
	"github.com/danyouknowme/smthng/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	hub := ws.NewWebsocketHub(&ws.Config{
		Redis: nil,
	})
	go hub.Run()

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to initialize server: %v", err)
		}
	}()

	logger.Infof("Listening on port %v", server.Addr)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources

	logger.Infof("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}
}
