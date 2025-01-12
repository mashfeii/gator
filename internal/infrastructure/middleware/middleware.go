package middleware

import (
	"context"
	"fmt"

	"github.com/mashfeii/gator/internal/domain"
)

type Middleware func(domain.Handler) domain.Handler

func Default(handler domain.Handler) domain.Handler {
	return handler
}

func LoggedIn(handler domain.Handler) domain.Handler {
	return func(s *domain.State, c *domain.Command) error {
		_, err := s.GetQueries().GetUser(context.Background(), s.GetConfig().CurrentUserName)
		if err != nil {
			return fmt.Errorf("user is not logged in: %w", err)
		}

		return handler(s, c)
	}
}
