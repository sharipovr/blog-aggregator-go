package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) != 1 {
		return errors.New("agg expects single argument, the time_between_reqs")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			log.Println(err)
		}
	}
}

func scrapeFeeds(s *state) error {

	dbFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	_, err = s.db.MarkFeedFetched(context.Background(), dbFeed.ID)
	if err != nil {
		return err
	}

	fmt.Println(dbFeed.Url)
	feed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return err
	}

	for _, feedItem := range feed.Channel.Item {
		fmt.Printf("%s | %s", feedItem.Title, feedItem.Link)
	}

	return nil
}
