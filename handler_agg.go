package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sharipovr/blog-aggregator-go/internal/database"
)

func handlerBrowse(s *state, cmd command) error {
	limit := 2
	var err error
	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post)
	}
	return nil
}

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
		var publishedAt sql.NullTime
		if feedItem.PubDate != "" {
			parsed, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
			if err == nil {
				publishedAt = sql.NullTime{Time: parsed, Valid: true}
			}
		}

		now := time.Now().UTC()
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title:     feedItem.Title,
			Url:       feedItem.Link,
			Description: sql.NullString{
				String: feedItem.Description,
				Valid:  feedItem.Description != "",
			},
			PublishedAt: publishedAt,
			FeedID:      dbFeed.ID,
		})

		if err != nil && !strings.Contains(err.Error(), "URL") {

			log.Println(err)
		}
	}
	return nil
}
