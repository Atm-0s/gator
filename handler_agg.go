package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Atm-0s/BlogAggregator/internal/database"
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
		fmt.Printf("-- %v\n", item.Title)
	}
	return nil
}
