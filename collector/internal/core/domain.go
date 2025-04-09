package core

import (
	"context"
	"time"
)

type OrderBookUpdate struct {
	Data      OrderBookData `bson:"data"`
	Timestamp time.Time     `bson:"timestamp"`
}

type OrderBookData struct {
	EventType     string     `json:"e"` // "depthUpdate"
	EventTime     int64      `json:"E"` // Event timestamp
	Symbol        string     `json:"s"` // "BTCUSDT"
	FirstUpdateID int64      `json:"U"` // First update ID in event
	FinalUpdateID int64      `json:"u"` // Final update ID in event
	BidUpdates    [][]string `json:"b"` // [["Price", "Quantity"],...]
	AskUpdates    [][]string `json:"a"` // [["Price", "Quantity"],...]
}

type OrderBookRepository interface {
	Save(ctx context.Context, update OrderBookUpdate) error
}

type WebSocketClient interface {
	Connect() error
	ReadMessages() (<-chan []byte, <-chan error)
	Close() error
}
