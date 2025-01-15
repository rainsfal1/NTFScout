package db

import (
	"context"
	"log"
	"nftscout/internal/api"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	client     *mongo.Client
	database   *mongo.Database
	collections *mongo.Collection
	transactions *mongo.Collection
	errors     *mongo.Collection
}

type ErrorLog struct {
	ID        string    `bson:"_id,omitempty"`
	Type      string    `bson:"type"`
	Message   string    `bson:"message"`
	Context   string    `bson:"context"`
	Timestamp time.Time `bson:"timestamp"`
}

func ConnectToMongoDB(ctx context.Context, uri, dbName string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}

func NewMongoDBPersister(uri, dbName string) (*MongoDBPersister, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	database := client.Database(dbName)
	
	return &MongoDBPersister{
		client:       client,
		database:     database,
		collections:  database.Collection("collections"),
		transactions: database.Collection("transactions"),
		errors:       database.Collection("errors"),
	}, nil
}

func (m *MongoDBPersister) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

func (client *MongoDBPersister) GetTransactionsFromDb(ctx context.Context, address string) (bool, error) {
	coll := client.database.Collection(os.Getenv("MONGODB_COLLECTION_TRANSACTION"))

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
	coll := client.database.Collection(os.Getenv("MONGODB_COLLECTION_TRANSACTION"))

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

func (m *MongoDBPersister) StoreCollection(collection api.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Add timestamp if not present
	if collection.Timestamp.IsZero() {
		collection.Timestamp = time.Now()
	}

	// Use upsert to avoid duplicates based on contract address
	filter := bson.M{"contract_address": collection.ContractAddress}
	update := bson.M{"$set": collection}
	opts := options.Update().SetUpsert(true)

	result, err := m.collections.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.UpsertedCount > 0 {
		log.Printf("Inserted new collection: %s", collection.Name)
	} else {
		log.Printf("Updated existing collection: %s", collection.Name)
	}

	return nil
}

func (m *MongoDBPersister) StoreTransaction(transaction api.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Add timestamp if not present
	if transaction.Timestamp.IsZero() {
		transaction.Timestamp = time.Now()
	}

	_, err := m.transactions.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}

	log.Printf("Stored transaction: %s", transaction.Hash)
	return nil
}

func (m *MongoDBPersister) LogError(errorType, message, context string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorLog := ErrorLog{
		Type:      errorType,
		Message:   message,
		Context:   context,
		Timestamp: time.Now(),
	}

	_, err := m.errors.InsertOne(ctx, errorLog)
	if err != nil {
		return err
	}

	log.Printf("Logged error [%s]: %s", errorType, message)
	return nil
}

func (m *MongoDBPersister) GetCollections(limit int) ([]api.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := m.collections.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var collections []api.Collection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, err
	}

	return collections, nil
}

func (m *MongoDBPersister) GetTransactions(limit int) ([]api.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := m.transactions.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []api.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
