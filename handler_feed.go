package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Atm-0s/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("addfeed expects a name and a url")
	}
	if len(cmd.Args) > 2 {
		return errors.New("too many arguments provided")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()
	user := s.cfgPtr.CurrentUserName
	userDB, err := s.db.GetUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	userID := userDB.ID

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    userID,
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}
	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)
	return nil
}
