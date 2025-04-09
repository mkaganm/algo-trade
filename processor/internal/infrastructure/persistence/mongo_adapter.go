package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/mkaganm/algo-trade/processor/internal/config"
	"github.com/mkaganm/algo-trade/processor/internal/core/domain"
	_ "github.com/mkaganm/algo-trade/processor/internal/core/ports/secondary" // fixme
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// fixme env loading

type MongoOrderBookRepository struct {
	client       *mongo.Client
	databaseName string
	collection   string
}

func NewMongoOrderBookRepository(uri, dbName, collection string) (*MongoOrderBookRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return &MongoOrderBookRepository{
		client:       client,
		databaseName: dbName,
		collection:   collection,
	}, nil
}

func (r *MongoOrderBookRepository) GetLatestRecords(ctx context.Context, limit int) ([]domain.OrderBookRecord, error) {
	r.collection = config.Load().CollectionName

	collection := r.client.Database(r.databaseName).Collection(r.collection)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "timestamp", Value: -1}})
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch records: %w", err)
	}
	defer cursor.Close(ctx)

	var records []domain.OrderBookRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("failed to decode records: %w", err)
	}

	return records, nil
}

func (r *MongoOrderBookRepository) SaveSignal(ctx context.Context, signal domain.TradeSignal) error {
	r.collection = config.Load().SignalsColName

	collection := r.client.Database(r.databaseName).Collection(r.collection)

	_, err := collection.InsertOne(ctx, signal)
	if err != nil {
		return fmt.Errorf("failed to insert signal: %w", err)
	}

	return nil
}
