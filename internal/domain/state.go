package domain

import (
	"github.com/mashfeii/gator/internal/config"
	"github.com/mashfeii/gator/internal/infrastructure/database"
)

type State struct {
	Queries *database.Queries
	Conf    *config.Config
}
