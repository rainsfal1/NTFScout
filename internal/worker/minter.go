package worker

import (
	"context"
	"log"
	"nftscout/internal/api"
	"nftscout/internal/db"
	"nftscout/internal/web3"
	"time"
)

func MintWorker(ctx context.Context, taskChan <-chan api.Collection, wallet *web3.Wallet, db *db.MongoDBPersister) {
	log.Println("Starting NFT minter worker...")

	for {
		select {
		case <-ctx.Done():
			log.Println("NFT minter stopped")
			return
		case collection := <-taskChan:
			log.Printf("Processing mint request for collection: %s", collection.Name)
			
			// Create transaction for minting
			tx, err := wallet.CreateTransaction(
				collection.ContractAddress,
				collection.Price,
				[]byte{}, // Empty data for basic transfer
			)
			if err != nil {
				log.Printf("Error creating transaction for %s: %v", collection.Name, err)
				
				// Log error to database
				if dbErr := db.LogError("mint_transaction_create", err.Error(), collection.ContractAddress); dbErr != nil {
					log.Printf("Error logging to database: %v", dbErr)
				}
				continue
			}

			// Send transaction
			txHash, err := wallet.SendTransaction(tx)
			if err != nil {
				log.Printf("Error sending transaction for %s: %v", collection.Name, err)
				
				// Log error to database
				if dbErr := db.LogError("mint_transaction_send", err.Error(), collection.ContractAddress); dbErr != nil {
					log.Printf("Error logging to database: %v", dbErr)
				}
				continue
			}

			// Store successful transaction
			transaction := api.Transaction{
				Hash:            txHash,
				From:            wallet.GetAddress(),
				To:              collection.ContractAddress,
				Value:           collection.Price,
				GasUsed:         0, // Will be updated when transaction is confirmed
				Status:          "pending",
				Timestamp:       time.Now(),
				CollectionName:  collection.Name,
			}

			if err := db.StoreTransaction(transaction); err != nil {
				log.Printf("Error storing transaction %s: %v", txHash, err)
			} else {
				log.Printf("Successfully initiated mint transaction %s for collection %s", txHash, collection.Name)
			}

			// Add delay to prevent rate limiting
			time.Sleep(2 * time.Second)
		}
	}
}
