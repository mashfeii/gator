package domain

import (
	"github.com/mashfeii/gator/internal/config"
	"github.com/mashfeii/gator/internal/infrastructure/database"
)

type State struct {
	queries *database.Queries
	conf    *config.Config
	user    *database.User
}

func (s State) GetQueries() *database.Queries {
	return s.queries
}

func (s *State) SetQueries(queries *database.Queries) {
	s.queries = queries
}

func (s State) GetConfig() *config.Config {
	return s.conf
}

func (s *State) SetConfig(conf *config.Config) {
	s.conf = conf
}

func (s State) GetUser() *database.User {
	return s.user
}

func (s *State) SetUser(user *database.User) {
	s.user = user
}
