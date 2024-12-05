package worker

import (
	"context"
	"log"
	"nftscout/internal/api"
	"time"
)

func FetchWorker(ctx context.Context, dataChan chan<- []api.Collection) {
	ticker := time.NewTicker(60 * time.Second) // Default 60 seconds
	defer ticker.Stop()

	log.Println("Starting NFT collection fetch worker...")

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
