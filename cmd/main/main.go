package main

import (
	"log/slog"

	"github.com/Vlad-Ali/Movies-service-back/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		slog.Error("Error creating application. ", "error", err)
		return
	}
	err = application.Start()
	if err != nil {
		slog.Error("Error starting or closing application. ", "error", err)
	}
}
