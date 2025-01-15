package main

import (
	"context"
	"log"
	"nftscout/internal/api"
	"nftscout/internal/db"
	"nftscout/internal/helpers"
	"nftscout/internal/web3"
	"nftscout/internal/worker"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment variables", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := db.ConnectToMongoDB(ctx, os.Getenv("MONGO_URI"), os.Getenv("MONGODB_DATABASE"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	persister := db.NewMongoDBPersister(client, os.Getenv("MONGODB_DATABASE"))

	gasLimit, err := strconv.Atoi(os.Getenv("GAS_LIMIT"))
	if err != nil {
		log.Fatal("Error converting GAS_LIMIT to integer", err)
	}
	
	wallet, err := web3.NewWallet(os.Getenv("PRIVATE_KEY"), os.Getenv("RPC_URL"), uint64(gasLimit))
	if err != nil {
		log.Fatal("Error creating wallet", err)
	}

	dataChannel := make(chan []api.Collection, 4)
	defer close(dataChannel)

	txChannel := make(chan helpers.ProcessedData, 4)
	defer close(txChannel)

	log.Println("NFTScout Bot Started")

	go worker.FetchWorker(ctx, dataChannel)
	go worker.TaskProcessor(ctx, dataChannel, txChannel)
	go worker.Minter(ctx, persister, txChannel, wallet)

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
