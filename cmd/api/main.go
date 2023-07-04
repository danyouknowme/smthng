package main

import (
	"github.com/danyouknowme/smthng/cmd/api/server"
	"github.com/danyouknowme/smthng/internal/config"
	"github.com/danyouknowme/smthng/pkg/logger"
)

var appConfig config.AppConfig
var err error

func init() {
	err = config.Load(&appConfig)
	if err != nil {
		logger.Fatal(err)
	}
}

func main() {
	app := server.NewApp(&appConfig)
	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
}
