package main

import (
	"context"
	"fmt"
)

func handlerFetchFeed(s *state, cmd command) error {
	ctx := context.Background()

	/*
		err := argCheck(
			cmd,
			"agg expects a valid URL",
		)
		if err != nil {
			return err
		}
	*/

	url := /* cmd.Args[0] */ "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
