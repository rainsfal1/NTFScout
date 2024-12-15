package web3

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	publicKey  *ecdsa.PublicKey
	address    common.Address
	gasLimit   uint64
}

func NewWallet(privateKeyHex, rpcURL string, gasLimit uint64) (*Wallet, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &Wallet{
		client:     client,
		privateKey: privateKey,
		publicKey:  publicKeyECDSA,
		address:    address,
		gasLimit:   gasLimit,
	}, nil
}

func (w *Wallet) GetAddress() common.Address {
	return w.address
}

func (w *Wallet) GetNonce(ctx context.Context) (uint64, error) {
	nonce, err := w.client.PendingNonceAt(ctx, w.address)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (w *Wallet) EstimateGasPrice(ctx context.Context) (*big.Int, error) {
	return w.client.SuggestGasPrice(ctx)
}

func (w *Wallet) CreateTransaction(ctx context.Context, to common.Address, value *big.Int, gasPrice *big.Int, data []byte) (*types.Transaction, error) {
	nonce, err := w.client.PendingNonceAt(ctx, w.address)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nonce, to, value, w.gasLimit, gasPrice, data)
	return tx, nil
}

func (w *Wallet) SendTransaction(ctx context.Context, tx *types.Transaction) (string, error) {
	chainID, err := w.client.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), w.privateKey)
	if err != nil {
		return "", err
	}

	err = w.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
