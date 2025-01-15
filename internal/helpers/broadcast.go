package helpers

import (
	"context"
	"log"
	"nftscout/internal/api"
)

func BroadCast(ctx context.Context, colls []api.Collection, txChan chan<- ProcessedData) {
	for _, col := range colls {
		// Check if context is canceled before processing
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping broadcast")
			return
		default:
		}

		transactions, err := api.GetTransaction(ctx, col)
		if err != nil {
			log.Println("Error getting transaction")
			continue
		}
		data, err := ProcessData(ctx, col, transactions)
		if err != nil {
			log.Println("Error processing transaction data")
			continue
		}

		// Safe channel send with context check
		select {
		case txChan <- data:
			// Successfully sent
		case <-ctx.Done():
			log.Println("Context canceled, stopping broadcast")
			return
		}
	}
}
