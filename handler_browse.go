package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jacobovii/gator/internal/database"
)

func handleBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.ListPostsForUser(context.Background(), database.ListPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.ListPostsForUserRow) {
	fmt.Printf("%s\n", post.PublishedAt.Time.Format("Mon Jan 2"))
	fmt.Printf("--- %s ---\n", post.Title)
	fmt.Printf("    %v\n", post.Description)
	fmt.Printf("Link: %s\n", post.Url)
	fmt.Println("=====================================")
}
