package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type TransactionData struct {
	Name            string
	Address         string
	Quantity        int64
	TransactionHash string
}

type DB interface {
	GetTransactions(ctx context.Context, address string) (bool, error)
	InsertTransaction(ctx context.Context, tx TransactionData) error
	LogError(ctx context.Context, err error)
}

type MongoDBPersister struct {
	MongoClient *mongo.Client
}

func NewMongoDBPersister(client *mongo.Client) *MongoDBPersister {
	return &MongoDBPersister{
		MongoClient: client,
	}
}

func ConnectToMongoDB(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return client, nil
}

func (client *MongoDBPersister) GetTransactionsFromDb(ctx context.Context, address string) (bool, error) {
	coll := client.MongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION_TRANSACTION"))

	filter := bson.D{{Key: "address", Value: address}}
	result := coll.FindOne(ctx, filter)

	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (client *MongoDBPersister) InsertTransactionToDb(ctx context.Context, tx TransactionData) error {
	coll := client.MongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION_TRANSACTION"))

	document := bson.M{
		"name":    tx.Name,
		"hash":    tx.TransactionHash,
		"qty":     tx.Quantity,
		"address": tx.Address,
	}

	if _, err := coll.InsertOne(ctx, document); err != nil {
		return err
	}
	return nil
}

func (client *MongoDBPersister) LogError(ctx context.Context, err error) {
	coll := client.MongoClient.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION_ERROR"))

	document := bson.M{
		"error":     err.Error(),
		"timeStamp": time.Now(),
	}

	coll.InsertOne(ctx, document)
}
