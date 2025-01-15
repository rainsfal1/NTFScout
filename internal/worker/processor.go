package worker

import (
	"context"
	"log"
	"nftscout/internal/api"
	"nftscout/internal/db"
)

func ProcessWorker(ctx context.Context, dataChan <-chan []api.Collection, taskChan chan<- api.Collection, db *db.MongoDBPersister) {
	log.Println("Starting data processor worker...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Data processor stopped")
			return
		case data := <-dataChan:
			log.Printf("Processing %d collections...", len(data))
			
			for _, collection := range data {
				// Store collection in database
				if err := db.StoreCollection(collection); err != nil {
					log.Printf("Error storing collection %s: %v", collection.Name, err)
					continue
				}
				
				// Send to task channel for minting
				select {
				case taskChan <- collection:
					log.Printf("Collection %s queued for minting", collection.Name)
				case <-ctx.Done():
					log.Println("Task channel closed, stopping processor")
					return
				}
			}
			
			log.Printf("Finished processing %d collections", len(data))
		}
	}
}
