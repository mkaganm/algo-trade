package domain

import "time"

const (
	Buy     = "BUY"
	Sell    = "SELL"
	Neutral = "NEUTRAL"
)

type TradeSignal struct {
	Signal    string    `bson:"signal"    json:"signal"`
	ShortSMA  float64   `bson:"shortSMA"  json:"shortSMA"`
	LongSMA   float64   `bson:"longSMA"   json:"longSMA"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}
