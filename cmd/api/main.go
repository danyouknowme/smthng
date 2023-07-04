package main

import (
	"github.com/danyouknowme/smthng/cmd/api/server"
	"github.com/danyouknowme/smthng/pkg/logger"
)

func main() {
	app := server.NewApp()
	if err := app.Start(); err != nil {
		logger.Fatal(err)
	}
}
