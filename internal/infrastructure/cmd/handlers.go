package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/mashfeii/gator/internal/domain"
	"github.com/mashfeii/gator/internal/infrastructure/database"
	"github.com/mashfeii/gator/internal/infrastructure/rss"
)

func LoginHandle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	user, err := s.GetQueries().GetUser(context.Background(), c.Args[0])
	if err != nil {
		return fmt.Errorf("user is not found: %w", err)
	}

	err = s.GetConfig().SetUser(c.Args[0])
	if err != nil {
		return fmt.Errorf("failed to switch to user: %w", err)
	}

	s.SetUser(&user)

	slog.Info("Login successful", "username", c.Args[0])

	return nil
}

func RegisterUserHandle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	user, err := s.GetQueries().CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to create a user: %w", err)
	}

	err = s.GetConfig().SetUser(user.Name)
	if err != nil {
		return err
	}

	s.SetUser(&user)

	slog.Info("User ID", "id", user.ID, "username", user.Name, "createdAt", user.CreatedAt, "updatedAt", user.UpdatedAt)

	return nil
}

func ResetHandle(s *domain.State, _ *domain.Command) error {
	err := s.GetQueries().Clear(context.Background())
	if err != nil {
		return err
	}

	err = s.GetConfig().SetUser("")
	if err != nil {
		return err
	}

	slog.Info("Database cleared")

	return nil
}

func GetUsersHandle(s *domain.State, _ *domain.Command) error {
	users, err := s.GetQueries().GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}

	for _, user := range users {
		fmt.Printf("* %s", user.Name)

		if s.GetConfig().CurrentUserName == user.Name {
			fmt.Printf(" (current)")
		}

		fmt.Println()
	}

	return nil
}

func CreateFeedHandle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	feed, err := s.GetQueries().CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
		Url:       c.Args[1],
		UserID:    s.GetUser().ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create a feed: %w", err)
	}

	_, err = s.GetQueries().CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    s.GetUser().ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create a follow to a feed: %w", err)
	}

	fmt.Printf("Created feed (%v):\n%v\t%v\t%v\n", s.GetUser().Name, feed.ID, feed.Name, feed.Url)

	return nil
}

func GetFeedsHandle(s *domain.State, _ *domain.Command) error {
	feeds, err := s.GetQueries().GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	for feed := range feeds {
		fmt.Printf("* %v\t%v\t%v\n", feeds[feed].Name, feeds[feed].Url, feeds[feed].UserID)
	}

	return nil
}

func FollowHandle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	feed, err := s.GetQueries().GetFeedURL(context.Background(), c.Args[0])
	if err != nil {
		return fmt.Errorf("unable to get feed, try to register first: %w", err)
	}

	row, err := s.GetQueries().CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    s.GetUser().ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create a feed follow: %w", err)
	}

	fmt.Printf("Successfully followed feed:\n%v\t%v\n", row.FeedsName, row.UsersName)

	return nil
}

func FollowingHandle(s *domain.State, _ *domain.Command) error {
	feeds, err := s.GetQueries().GetFeedFollowsForUser(context.Background(), s.GetUser().ID)
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	for i := 0; i != len(feeds); i++ {
		fmt.Printf("* %v\t%v\n", feeds[i].FeedsName, feeds[i].Url)
	}

	return nil
}

func Unfollow(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	err := s.GetQueries().DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: s.GetUser().ID,
		Url:    c.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to unfollow the feed: %w", err)
	}

	err = s.GetQueries().UpdateFeedUser(context.Background(), database.UpdateFeedUserParams{
		UserID: s.GetUser().ID,
		Url:    c.Args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to unfollow the feed: %w", err)
	}

	fmt.Printf("Successfully unfollowed:\n%s\t%s", s.GetUser().Name, c.Args[0])

	return nil
}

func UpdateFeedsHandle(s *domain.State, c *domain.Command) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	duration, err := time.ParseDuration(c.Args[0])
	if err != nil {
		return fmt.Errorf("unable to get time duration: %w", err)
	}

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		err := updateFeeds(s, c)
		if err != nil {
			slog.Error("error on feed fetching", "error", err.Error())
		}
	}
}

func updateFeeds(s *domain.State, _ *domain.Command) error {
	nextFeed, err := s.GetQueries().GetNextFeedToFetch(context.Background(), s.GetUser().ID)
	if err != nil {
		return fmt.Errorf("unable to get next feed to update: %w", err)
	}

	err = s.GetQueries().MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        nextFeed.ID,
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("unable to update feed: %w", err)
	}

	updateFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("unable to get updated feed: %w", err)
	}

	validItems := 0

	for item := range updateFeed.Channel.Item {
		published, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", updateFeed.Channel.Item[item].PubDate)
		if err != nil {
			slog.Error("unable to parse date", "error", err.Error())
			continue
		}

		descr, descrValid := updateFeed.Channel.Item[item].Description, true
		if descr == "" {
			descrValid = false
		}

		_, err = s.GetQueries().CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     updateFeed.Channel.Item[item].Title,
			Description: sql.NullString{
				String: descr,
				Valid:  descrValid,
			},
			Url:         updateFeed.Channel.Item[item].Link,
			PublishedAt: published,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				slog.Error("unable to create post", "error", err.Error())
			}

			continue
		}

		validItems++
	}

	slog.Info("Feeds updated", "number of new posts", validItems, "feed", nextFeed.Name)

	return nil
}

func Browse(s *domain.State, c *domain.Command) error {
	var (
		limit    = 2
		limitErr error
	)

	if len(c.Args) == 1 {
		limit, limitErr = strconv.Atoi(c.Args[0])
		if limitErr != nil {
			return fmt.Errorf("unable to convert limit value: %w", limitErr)
		}
	}

	posts, err := s.GetQueries().GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: s.GetUser().ID,
		Limit:  int32(limit), //nolint // no possible overflow
	})
	if err != nil {
		return fmt.Errorf("unable to get feeds: %w", err)
	}

	fmt.Printf("- Posts for %s\n", s.GetUser().Name)

	for post := range posts {
		postTitle, postLink, postDate := posts[post].Title, posts[post].Url, posts[post].PublishedAt

		fmt.Printf("* %v:\n\tPost title: %v\n\tPost link: %v\n\tPost date: %v\n", posts[post].FeedName, postTitle, postLink, postDate)
	}

	return nil
}
