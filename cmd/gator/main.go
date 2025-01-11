package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"github.com/mashfeii/gator/internal/application"
	"github.com/mashfeii/gator/internal/domain"
)

func main() {
	if len(os.Args) == 0 {
		slog.Error("No command provided")
		os.Exit(1)
	}

	app, err := application.NewApp()
	if err != nil {
		slog.Error("Failed to initialize app", "error", err)
		os.Exit(1)
	}

	err = app.Commands.Run(&app.State, &domain.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	})
	if err != nil {
		slog.Error("Failed to run command", "error", err)
		os.Exit(1)
	}
}
