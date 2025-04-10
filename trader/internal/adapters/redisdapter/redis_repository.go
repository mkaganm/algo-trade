package redisdapter

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) ReadMessages(ctx context.Context) ([]map[string]interface{}, error) {
	// Ensure the consumer group exists
	err := r.client.XGroupCreateMkStream(ctx, "trade_signals_stream", "signal_group", "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Read messages from the stream
	streams, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "signal_group",
		Consumer: "signal_consumer",
		Streams:  []string{"trade_signals_stream", ">"},
		Count:    10,
		Block:    0,
	}).Result()
	if err != nil {
		return nil, err
	}

	var messages []map[string]interface{}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			message.Values["id"] = message.ID
			messages = append(messages, message.Values)
		}
	}

	return messages, nil
}

// CheckHealth implements the HealthService interface.
func (r *RedisRepository) CheckHealth(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()

	return err
}

func (r *RedisRepository) WriteProcessedMessage(ctx context.Context, message map[string]interface{}) error {
	_, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "processed_signals",
		Values: message,
	}).Result()

	return err
}

func (r *RedisRepository) AcknowledgeMessage(ctx context.Context, messageID string) error {
	return r.client.XAck(ctx, "trade_signals_stream", "signal_group", messageID).Err()
}
