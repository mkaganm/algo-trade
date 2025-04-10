package ports

import (
	"context"

	"github.com/mkaganm/algo-trade/processor/internal/core/domain"
)

// OrderBookRepository is the secondary port (interface) for order book data access.
type OrderBookRepository interface {
	GetLatestRecords(ctx context.Context, limit int) ([]domain.OrderBookRecord, error)
}

// SignalRepository is the secondary port (interface) for signal storage.
type SignalRepository interface {
	SaveSignal(ctx context.Context, signal domain.TradeSignal) error
}

// SignalPublisher is the secondary port (interface) for publishing signals.
type SignalPublisher interface {
	PublishSignal(ctx context.Context, signal domain.TradeSignal) error
}
