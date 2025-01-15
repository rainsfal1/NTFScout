package helpers

import (
	"context"
	"log"
	"math/big"
	"nftscout/internal/db"
	"nftscout/internal/web3"

	"github.com/ethereum/go-ethereum/common"
)

func Transaction(ctx context.Context, data ProcessedData, persister *db.MongoDBPersister, wallet *web3.Wallet) {
	isMinted, err := persister.GetTransactionsFromDb(ctx, data.Contract)
	if err != nil {
		return
	}

	if isMinted {
		log.Println("Collection already minted")
		return
	}

	gasPrice, err := wallet.EstimateGasPrice(ctx)
	if err != nil {
		return
	}

	tx, err := wallet.CreateTransaction(ctx, common.HexToAddress(data.To), big.NewInt(0), gasPrice, []byte(data.CallData))

	if err != nil {
		return
	}

	hash, err := wallet.SendTransaction(ctx, tx)
	if err != nil {
		return
	}

	txData := db.TransactionData{
		Name:            data.Name,
		Address:         data.Contract,
		Quantity:        data.Quantity,
		TransactionHash: hash,
	}

	err = persister.InsertTransactionToDb(ctx, txData)
	if err != nil {
		return
	}

}
