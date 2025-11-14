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
	c := command{
		Args: []string{url},
	}
	err = handlerFollowFeed(s, c)
	if err != nil {
		return err
	}

	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("feeds command does not take arguments")
	}
	ctx := context.Background()
	feedSlice, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error showing feeds")
	}
	for _, feed := range feedSlice {
		fmt.Println("===================")
		fmt.Printf("Feed: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		fmt.Printf("User: %v\n", feed.User)
		fmt.Println("===================")
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	err := argCheck(cmd, "follow requires a url")
	if err != nil {
		return err
	}

	ctx := context.Background()
	url := cmd.Args[0]

	userDB, err := s.db.GetUser(ctx, s.cfgPtr.CurrentUserName)
	if err != nil {
		return err
	}
	feedDB, err := s.db.GetFeedFromURL(ctx, url)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userDB.ID,
		FeedID:    feedDB.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, followParams)
	if err != nil {
		return fmt.Errorf("error registering follow: %w", err)
	}
	fmt.Printf("User: %s", userDB.Name)
	fmt.Println("is now following")
	fmt.Printf("Feed: %v", feedDB.Name)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("following does not take any arguments")
	}
	ctx := context.Background()
	currentUser, err := s.db.GetUser(ctx, s.cfgPtr.CurrentUserName)
	if err != nil {
		return err
	}
	userFeeds, err := s.db.GetFeedFollowsForUser(ctx, currentUser.ID)
	if err != nil {
		return fmt.Errorf("error fetching feeds for user %s: %w", currentUser.Name, err)
	}

	fmt.Printf("User %s is following:\n", currentUser.Name)

	for _, feed := range userFeeds {
		f, err := s.db.GetFeedByID(ctx, feed.FeedID)
		if err != nil {
			return err
		}
		fmt.Println(f.Name)
	}
	return nil
}
