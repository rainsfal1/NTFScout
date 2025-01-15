package worker

import (
	"context"
	"log"
	"nftscout/internal/api"
	"os"
	"strconv"
	"time"
)

func FetchWorker(ctx context.Context, dataChan chan<- []api.Collection) {
	// Get fetch duration from environment variable
	fetchDurationStr := os.Getenv("FETCH_DURATION")
	if fetchDurationStr == "" {
		fetchDurationStr = "60" // default to 60 seconds
	}
	
	fetchDuration, err := strconv.Atoi(fetchDurationStr)
	if err != nil {
		log.Printf("Invalid FETCH_DURATION: %v, using default 60 seconds", err)
		fetchDuration = 60
	}
	
	ticker := time.NewTicker(time.Duration(fetchDuration) * time.Second)
	defer ticker.Stop()

	log.Printf("Starting fetch worker with %d second intervals", fetchDuration)

	for {
		select {
		case <-ctx.Done():
			log.Println("Collection fetcher stopped")
			return
		case <-ticker.C:
			log.Println("Fetching collection data...")
			data, err := api.FetchCollection(ctx)
			if err != nil {
				log.Printf("Error fetching collection data: %v", err)
				continue
			}

			log.Printf("Successfully fetched %d collections", len(data))
			select {
			case dataChan <- data:
				log.Println("Data sent to processing channel")
			case <-ctx.Done():
				log.Println("Data channel closed, stopping fetcher")
				return
			}
		}
	}
}
