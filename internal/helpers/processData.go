package helpers

import (
	"context"
	"fmt"
	"nftscout/internal/api"
	"strconv"
)

type ProcessedData struct {
	Name     string
	Contract string
	Quantity int64
	EthValue string
	To       string
	CallData string
	Hash     string
}

func ProcessData(ctx context.Context, collection api.Collection, txs []api.Transaction) (ProcessedData, error) {
	if len(txs) == 0 {
		return ProcessedData{}, fmt.Errorf("no items provided")
	}

	var minTransaction api.Transaction
	minNftCount := int(^uint(0) >> 1)

	for _, tx := range txs {
		nftCount, err := strconv.Atoi(tx.NftCount)
		if err != nil {
			return ProcessedData{}, fmt.Errorf("invalid NFT count: %w", err)
		}
		if nftCount < minNftCount {
			minNftCount = nftCount
			minTransaction = tx
		}
	}

	minNftCount, err := strconv.Atoi(minTransaction.NftCount)
	if err != nil {
		return ProcessedData{}, fmt.Errorf("invalid NFT count: %w", err)
	}
	data := ProcessedData{
		Name:     collection.Name,
		Contract: collection.Contract,
		Quantity: int64(minNftCount),
		To:       minTransaction.To,
		CallData: minTransaction.CallData,
		EthValue: minTransaction.EthValue,
	}
	return data, nil
}
