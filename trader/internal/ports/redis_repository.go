package ports

import "context"

type RedisRepository interface {
	ReadMessages(ctx context.Context) ([]map[string]interface{}, error)
	WriteProcessedMessage(ctx context.Context, message map[string]interface{}) error
	AcknowledgeMessage(ctx context.Context, messageID string) error
}
