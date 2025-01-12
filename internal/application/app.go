package application

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mashfeii/gator/internal/config"
	"github.com/mashfeii/gator/internal/domain"
	"github.com/mashfeii/gator/internal/infrastructure/cmd"
	"github.com/mashfeii/gator/internal/infrastructure/database"
	"github.com/mashfeii/gator/internal/infrastructure/middleware"
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
	app.Commands.Register("login", cmd.LoginHandle)
	app.Commands.Register("register", cmd.RegisterUserHandle)
	app.Commands.Register("reset", cmd.ResetHandle)
	app.Commands.Register("users", cmd.GetUsersHandle)
	app.Commands.Register("agg", middleware.LoggedIn(cmd.UpdateFeedsHandle))
	app.Commands.Register("addfeed", middleware.LoggedIn(cmd.CreateFeedHandle))
	app.Commands.Register("feeds", cmd.GetFeedsHandle)
	app.Commands.Register("follow", middleware.LoggedIn(cmd.FollowHandle))
	app.Commands.Register("following", middleware.LoggedIn(cmd.FollowingHandle))
	app.Commands.Register("unfollow", middleware.LoggedIn(cmd.Unfollow))
	app.Commands.Register("browse", middleware.LoggedIn(cmd.Browse))

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

	// NOTE: set current user
	user, err := queries.GetUser(context.Background(), conf.CurrentUserName)
	if err != nil && conf.CurrentUserName != "" {
		user, err = queries.CreateUser(context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			Name:      conf.CurrentUserName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	app.State.SetConfig(conf)
	app.State.SetQueries(queries)
	app.State.SetUser(&user)

	return app, nil
}
