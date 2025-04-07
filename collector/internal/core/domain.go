package core

import (
	"context"
	"time"
)

type OrderBookUpdate struct {
	Data      string    `bson:"data"`
	Timestamp time.Time `bson:"timestamp"`
}

type OrderBookRepository interface {
	Save(ctx context.Context, update OrderBookUpdate) error
}

type WebSocketClient interface {
	Connect() error
	ReadMessages() (<-chan []byte, <-chan error)
	Close() error
}
