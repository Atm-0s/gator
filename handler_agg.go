package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Atm-0s/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("agg expects a single duration e.g 5s or 2m or 3h")
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	t_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	lowerLimit, err := time.ParseDuration("1s")
	if err != nil {
		return err
	}

	if t_between_reqs < lowerLimit {
		return fmt.Errorf("select a time >= 1s")
	}

	fmt.Printf("Collecting feeds every %v\n", t_between_reqs)

	ticker := time.NewTicker(t_between_reqs)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	fetchedParams := database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err = s.db.MarkFeedFetched(ctx, fetchedParams)
	if err != nil {
		return err
	}
	f, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}
	fmt.Println("========================")
	fmt.Printf("Feed: %v\n", feed.Name)
	fmt.Println("========================")

	for _, item := range f.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)

		postParams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time:  publishedAt,
				Valid: err == nil,
			},
			FeedID: feed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				continue
			} else {
				return fmt.Errorf("error creating posts: %w", err)
			}
		}

	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	limit = 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments")
	} else if len(cmd.Args) == 1 {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
		limit = int32(n)
	}
	pUserParams := database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	posts, err := s.db.GetPostsByUser(context.Background(), pUserParams)
	if err != nil {
		return err
	}
	fmt.Println("Browsing posts")
	fmt.Println("=====================")
	for _, post := range posts {
		fmt.Println(post.Title)
		fmt.Println(post.PublishedAt)
		fmt.Println(post.Description)
		fmt.Println(post.Url)
		fmt.Println("=====================")
	}
	fmt.Println("End of posts")
	return nil
}
