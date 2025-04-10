package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mkaganm/algo-trade/trader/internal/ports"
)

type MessageProcessor struct {
	redisRepo ports.RedisRepository
}

func NewMessageProcessor(redisRepo ports.RedisRepository) *MessageProcessor {
	return &MessageProcessor{redisRepo: redisRepo}
}

func (mp *MessageProcessor) ProcessMessages(ctx context.Context) {
	messages, err := mp.redisRepo.ReadMessages(ctx)
	if err != nil {
		log.Printf("Error reading messages: %v", err)

		return
	}

	msg := messages[0]

	id, ok := msg["id"].(string)
	if !ok {
		log.Printf("Message missing 'id' or invalid type: %v", msg)

		return
	}

	timestamp, ok := msg["time"].(string)
	if !ok {
		log.Printf("Message missing 'time' or invalid type: %v", msg)

		return
	}

	// fixme : handle errors
	processedData, _ := json.Marshal(msg)

	err = mp.redisRepo.WriteProcessedMessage(ctx, map[string]interface{}{
		"original_id":   id,
		"original_data": processedData,
		"processed":     true,
		"processed_at":  timestamp,
	})
	if err != nil {
		log.Printf("Error writing processed message: %v", err)

		return
	}

	err = mp.redisRepo.AcknowledgeMessage(ctx, id)
	if err != nil {
		log.Printf("Error acknowledging message: %v", err)

		return
	}

	mp.tradeProcess(msg)
}

// tradeProcess processes the trade signal and executes the corresponding action.
// It can be extended to include more complex logic, such as placing orders or updating positions.
// For now, it simply logs the action taken based on the signal.
// You can implement the actual trading logic here, such as placing orders through an exchange API.
func (mp *MessageProcessor) tradeProcess(msg map[string]interface{}) {
	switch msg["signal"] {
	case "BUY":
		log.Println("Executing BUY order")
	case "SELL":
		log.Println("Executing SELL order")
	case "NEUTRAL":
		log.Println("Holding position")
	default:
		log.Println("Unknown signal received")
	}
}
