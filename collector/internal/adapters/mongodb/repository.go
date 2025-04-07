package mongodb

import (
	"context"
	"time"

	"github.com/mkaganm/algo-trade/collector/internal/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	closeTimeout  = 5 * time.Second
	cancelTimeout = 10 * time.Second
	expireTTL     = 7 * 24 * 60 * 60 // 7 days
)

type MongoOrderBookRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoOrderBookRepository(uri, database, collection string) (*MongoOrderBookRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	repo := &MongoOrderBookRepository{
		client:     client,
		database:   database,
		collection: collection,
	}

	// create collection with TTL index
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"createdAt": 1},
		Options: options.Index().SetExpireAfterSeconds(expireTTL), // TTL time
	}

	_, err = client.Database(database).Collection(collection).Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (m *MongoOrderBookRepository) Save(ctx context.Context, update core.OrderBookUpdate) error {
	coll := m.client.Database(m.database).Collection(m.collection)

	_, err := coll.InsertOne(ctx, bson.M{
		"data":       update.Data,
		"timestamp":  update.Timestamp,
		"created_at": time.Now(),
	})

	return err
}

func (m *MongoOrderBookRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()

	return m.client.Disconnect(ctx)
}
