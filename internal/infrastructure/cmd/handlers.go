package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/mashfeii/gator/internal/domain"
	"github.com/mashfeii/gator/internal/infrastructure/database"
	"github.com/mashfeii/gator/internal/infrastructure/rss"
)

type Login struct{}

func (l Login) Handle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	_, err := s.Queries.GetUser(context.Background(), c.Args[0])
	if err != nil {
		return fmt.Errorf("user is not found: %w", err)
	}

	err = s.Conf.SetUser(c.Args[0])
	if err != nil {
		return fmt.Errorf("failed to switch to user: %w", err)
	}

	slog.Info("Login successful", "username", c.Args[0])

	return nil
}

type RegisterUser struct{}

func (r RegisterUser) Handle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	user, err := s.Queries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to create a user: %w", err)
	}

	err = s.Conf.SetUser(user.Name)
	if err != nil {
		return err
	}

	slog.Info("User ID", "id", user.ID, "username", user.Name, "createdAt", user.CreatedAt, "updatedAt", user.UpdatedAt)

	return nil
}

type Reset struct{}

func (r Reset) Handle(s *domain.State, _ *domain.Command) error {
	err := s.Queries.Clear(context.Background())
	if err != nil {
		return err
	}

	slog.Info("Database cleared")

	return nil
}

type GetUsers struct{}

func (u GetUsers) Handle(s *domain.State, _ *domain.Command) error {
	users, err := s.Queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)

		if s.Conf.CurrentUsername == user.Name {
			fmt.Printf(" (current)")
		}

		fmt.Println()
	}

	return nil
}

type FetchFeed struct{}

func (f FetchFeed) Handle(_ *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	feed, err := rss.FetchFeed(context.Background(), c.Args[0])
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}

	fmt.Printf("Fetched structure:\n%v\t%v\t%v\n", feed.Channel.Title, feed.Channel.Description, feed.Channel.Link)

	return nil
}

type CreateFeed struct{}

func (cr CreateFeed) Handle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	user, err := s.Queries.GetUser(context.Background(), s.Conf.CurrentUsername)
	if err != nil {
		return fmt.Errorf("unable to get current user: %w", err)
	}

	feed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
		Url:       c.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create a feed: %w", err)
	}

	fmt.Printf("Created feed:\n%v\t%v\t%v\n", feed.ID, feed.Name, feed.Url)

	return nil
}
