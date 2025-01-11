package application

import (
	"database/sql"
	"fmt"

	"github.com/mashfeii/gator/internal/config"
	"github.com/mashfeii/gator/internal/domain"
	"github.com/mashfeii/gator/internal/infrastructure/cmd"
	"github.com/mashfeii/gator/internal/infrastructure/database"
)

type App struct {
	Commands domain.Commands
	State    domain.State
}

func NewApp() (*App, error) {
	// NOTE: initial app structure
	app := &App{
		Commands: domain.Commands{
			Handlers: map[string]domain.Handler{},
		},
		State: domain.State{},
	}

	// NOTE: register commands
	app.Commands.Register("login", cmd.Login{})
	app.Commands.Register("register", cmd.RegisterUser{})
	app.Commands.Register("reset", cmd.Reset{})
	app.Commands.Register("users", cmd.GetUsers{})
	app.Commands.Register("agg", cmd.FetchFeed{})
	app.Commands.Register("addfeed", cmd.CreateFeed{})

	// NOTE: open connection to database
	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// NOTE: extract queries to manipulate database
	queries := database.New(db)

	// NOTE: read config file
	conf, err := config.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	app.State.Conf = conf
	app.State.Queries = queries

	return app, nil
}
