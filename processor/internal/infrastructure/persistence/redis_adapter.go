package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mkaganm/algo-trade/processor/internal/core/domain"
	_ "github.com/mkaganm/algo-trade/processor/internal/core/ports/secondary" // fixme :
)

type RedisSignalPublisher struct {
	client    *redis.Client
	streamKey string
}

func NewRedisSignalPublisher(addr, streamKey string) *RedisSignalPublisher {
	return &RedisSignalPublisher{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		streamKey: streamKey,
	}
}

func (p *RedisSignalPublisher) PublishSignal(ctx context.Context, signal domain.TradeSignal) error {
	_, err := p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: p.streamKey,
		Values: map[string]interface{}{
			"signal":   signal.Signal,
			"shortSMA": signal.ShortSMA,
			"longSMA":  signal.LongSMA,
			"time":     signal.Timestamp.Format(time.RFC3339),
		},
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to add signal to Redis Stream: %w", err)
	}

	return nil
}
