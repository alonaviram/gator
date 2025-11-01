package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alonaviram/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	req.Header.Set("User-Agent", "gator")
	if err != nil {
		return &RSSFeed{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var rssFeed RSSFeed

	err = xml.Unmarshal(bytes, &rssFeed)
	if err != nil {
		return &RSSFeed{}, err
	}
	return &rssFeed, nil
}

func handlerAgg(s *state, cmd command) error {
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println("tick")
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	dbFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("fetching %v\n", dbFeed.Name)

	feed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return err
	}
	_, err = s.db.MarkFeedFetched(context.Background(), dbFeed.ID)
	if err != nil {
		return err
	}
	for _, rssItem := range feed.Channel.Item {

		t, err := time.Parse(time.RFC1123Z, rssItem.PubDate)
		if err != nil {
			println("there was an error parsing date")
			t = time.Now()
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     rssItem.Title,
			Url:       rssItem.Link,
			Description: sql.NullString{
				String: rssItem.Description,
				Valid:  rssItem.Description != "",
			},
			PublishedAt: t,
			FeedID:      dbFeed.ID,
		})
		if err != nil {
			isDup := strings.Contains(err.Error(), "posts_url_key")
			if !isDup {
				fmt.Printf("%v", err.Error())
			}
		}
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("not enough parameters")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Created New Feed:\n%v\n", feed)
	return nil
}

func handlerShowFeed(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println("------------------")
		fmt.Printf("Name: %v\n", feed.Name)
		fmt.Printf("URL: %v\n", feed.Url)
		u, _ := s.db.GetUserById(context.Background(), feed.UserID)
		fmt.Printf("User Name: %v\n", u.Name)
		fmt.Println("------------------")
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32

	if len(cmd.args) < 1 {
		limit = 10
	} else {
		limitArg := cmd.args[0]
		limitInt, err := strconv.Atoi(limitArg)
		if err != nil {
			return fmt.Errorf("limit arg must be of type int, got:%v", limitArg)
		}

		limit = int32(limitInt)
	}
	posts, error := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if error != nil {
		return error
	}
	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	// fmt.
	return nil

	// timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	// if err != nil {
	// 	return err
	// }
	// ticker := time.NewTicker(timeBetweenRequests)
	// for ; ; <-ticker.C {
	// 	fmt.Println("tick")
	// 	scrapeFeeds(s)
	// }
}
