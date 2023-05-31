package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/google/uuid"
)

func fetchFeed(url string) (*RSSFeed, error){
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	
	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	err = xml.Unmarshal(data, &feed)

	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func scrapeFeed(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	_, err := db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		log.Printf("Couldn't mark feed %s as fetched, Err: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't fetch feed %s, Err: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		log.Println("Found post: ", item.Title)

		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}

			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	

	log.Printf("Feed %s fetched, found %v posts", feed.Name, len(feedData.Channel.Item))
}

func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}

		log.Printf("Found %v feeds to fetch!", len(feeds))

		waitGroup := &sync.WaitGroup{}
		
		for _, feed := range feeds {
			waitGroup.Add(1)
			go scrapeFeed(db, waitGroup, feed)
		}
		waitGroup.Wait()
	}
}
