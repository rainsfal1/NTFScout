package web3

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Web3Wallet interface {
	GetNonce(ctx context.Context) (uint64, error)
	EstimateGasPrice(ctx context.Context) (*big.Int, error)
	CreateTransaction(ctx context.Context, toAddress common.Address, value *big.Int, gasPrice *big.Int, data []byte) (*types.Transaction, error)
	SendTransaction(ctx context.Context, signedTx *types.Transaction) error
}

type Wallet struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	gasLimit   uint64
	address    common.Address
}

func NewWallet(privateKey string, rpcUrl string, gasLimit uint64) (*Wallet, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	pKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		client:     client,
		privateKey: pKey,
		gasLimit:   gasLimit,
		address:    crypto.PubkeyToAddress(pKey.PublicKey),
	}

	return wallet, nil
}

func (w *Wallet) GetNonce(ctx context.Context) (uint64, error) {
	nonce, err := w.client.PendingNonceAt(ctx, w.address)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (w *Wallet) EstimateGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := w.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	return gasPrice, nil
}

func (w *Wallet) CreateTransaction(ctx context.Context, toAddress common.Address, value *big.Int, gasPrice *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := w.GetNonce(ctx)
	if err != nil {
		log.Println("Error getting nonce")
		return nil, err
	}

	tx := types.NewTransaction(nonce, toAddress, value, w.gasLimit, gasPrice, data)

	chainId, err := w.client.ChainID(ctx)
	if err != nil {
		log.Println("Error getting chain ID")
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), w.privateKey)
	if err != nil {
		log.Println("Error while signing transaction")
		return nil, err
	}

	return signedTx, nil
}

func (w *Wallet) SendTransaction(ctx context.Context, signedTx *types.Transaction) (string, error) {
	err := w.client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Println("Error while sending transaction")
		return "", err
	}

	txHash := signedTx.Hash().Hex()

	return txHash, nil
}
