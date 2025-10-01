package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alonaviram/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("expected url parameter")
	}
	url := cmd.args[0]

	f, e := s.db.GetFeedByUrl(context.Background(), url)
	if e != nil {
		return e
	}
	feedFollowRow, e := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    f.ID,
	})
	if e != nil {
		return e
	}
	fmt.Printf("name: %v, user: %v\n", feedFollowRow.FeedName, feedFollowRow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feedFollow := range feedFollows {
		fmt.Printf("%v\n", feedFollow.FeedName)
	}

	// f, e := s.db.GetFeedByUrl(context.Background(), url)
	// if e != nil {
	// 	return e
	// }
	// feedFollowRow, e := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
	// 	ID:        uuid.New(),
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// 	UserID:    u.ID,
	// 	FeedID:    f.ID,
	// })
	// if e != nil {
	// 	return e
	// }
	// fmt.Printf("name: %v, user: %v\n", feedFollowRow.FeedName, feedFollowRow.UserName)

	return nil
}
