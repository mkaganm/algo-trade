package domain

import "time"

type OrderBookData struct {
	EventType     string     `json:"e"`
	EventTime     int64      `json:"E"`
	Symbol        string     `json:"s"`
	FirstUpdateID int64      `json:"U"`
	FinalUpdateID int64      `json:"u"`
	BidUpdates    [][]string `json:"b"`
	AskUpdates    [][]string `json:"a"`
}

type OrderBookRecord struct {
	Data      OrderBookData `bson:"data"`
	Timestamp time.Time     `bson:"timestamp"`
}
