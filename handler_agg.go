package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jacobovii/gator/internal/database"
)

func handleGetRSSFeed(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenReq, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("please provide a valid duration (e.g. 1s, 1m, 1h, etc): %w", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReq)
	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Couldn't get next feed to fetch: %v", err)
		return
	}
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed fetched: %v", err)
	}

	fmt.Printf("Fetching new content for %s...\n", feed.Name)
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldn't collect feed: %v", err)
	}

	for _, item := range rssFeed.Channel.Item {
		log.Printf(" Saving post: %s\n", item.Title)

		publishedAt := sql.NullTime{}
		layout := "Mon, 02 Jan 2006 15:04:05 Z0700"
		if parsedTime, err := time.Parse(layout, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
