package worker

import (
	"context"
	"log"
	"nftscout/internal/db"
	"nftscout/internal/helpers"
	"nftscout/internal/web3"
)

func Minter(ctx context.Context, db *db.MongoDBPersister, txChan <-chan helpers.ProcessedData, wallet *web3.Wallet) {
	for {
		select {
		case <-ctx.Done():
			log.Println("NFTScout stopped")
			return
		case tx := <-txChan:
			helpers.Transaction(ctx, tx, db, wallet)
		}
	}
}
