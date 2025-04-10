package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mkaganm/algo-trade/processor/internal/core/domain"
	_ "github.com/mkaganm/algo-trade/processor/internal/core/ports/primary" // fixme
	"github.com/mkaganm/algo-trade/processor/internal/core/ports/secondary"
)

// Define static errors.
var (
	ErrNotEnoughDataPoints = errors.New("not enough data points to calculate SMAs")
	ErrNotEnoughData       = errors.New("not enough data to process signals")
)

type SignalProcessor struct {
	orderBookRepo secondary.OrderBookRepository
	signalRepo    secondary.SignalRepository
	publisher     secondary.SignalPublisher
}

func NewSignalProcessor(
	orderBookRepo secondary.OrderBookRepository,
	signalRepo secondary.SignalRepository,
	publisher secondary.SignalPublisher,
) *SignalProcessor {
	return &SignalProcessor{
		orderBookRepo: orderBookRepo,
		signalRepo:    signalRepo,
		publisher:     publisher,
	}
}

func (s *SignalProcessor) CalculateSMAs(
	prices []float64, shortPeriod,
	longPeriod int,
) (shortSMA, longSMA []float64, err error) {
	if len(prices) < longPeriod {
		return nil, nil, ErrNotEnoughDataPoints
	}

	shortSMA = calculateSMA(prices, shortPeriod)
	longSMA = calculateSMA(prices, longPeriod)

	return shortSMA, longSMA, nil
}

func (s *SignalProcessor) GenerateSignal(
	ctx context.Context,
	shortPeriod int,
	longPeriod int,
) (*domain.TradeSignal, error) {
	records, err := s.orderBookRepo.GetLatestRecords(ctx, longPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book records: %w", err)
	}

	if len(records) < longPeriod {
		return nil, ErrNotEnoughData
	}

	// Extract prices
	var prices []float64

	prices, err = extractPrices(records)
	if err != nil {
		return nil, fmt.Errorf("failed to extract prices: %w", err)
	}

	shortSMA, longSMA, err := s.CalculateSMAs(prices, shortPeriod, longPeriod)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate SMAs: %w", err)
	}

	if len(shortSMA) == 0 || len(longSMA) == 0 {
		return nil, ErrNotEnoughDataPoints
	}

	lastShortSMA := shortSMA[len(shortSMA)-1]
	lastLongSMA := longSMA[len(longSMA)-1]

	signal := selectSignal(lastShortSMA, lastLongSMA)

	tradeSignal := &domain.TradeSignal{
		Signal:    signal,
		ShortSMA:  lastShortSMA,
		LongSMA:   lastLongSMA,
		Timestamp: time.Now(),
	}

	// Save to database
	if err := s.signalRepo.SaveSignal(ctx, *tradeSignal); err != nil {
		log.Printf("Failed to save signal to database: %v", err)
	}

	// Publish to Redis
	if err := s.publisher.PublishSignal(ctx, *tradeSignal); err != nil {
		log.Printf("Failed to publish signal: %v", err)
	}

	return tradeSignal, nil
}

func selectSignal(lastShortSMA, lastLongSMA float64) string {
	var signal string

	switch {
	case lastShortSMA > lastLongSMA:
		signal = domain.Buy
	case lastShortSMA < lastLongSMA:
		signal = domain.Sell
	default:
		signal = domain.Neutral
	}

	return signal
}

func extractPrices(records []domain.OrderBookRecord) ([]float64, error) {
	var prices []float64

	for _, record := range records {
		if len(record.Data.BidUpdates) > 0 {
			price := record.Data.BidUpdates[0][0]

			var p float64

			_, err := fmt.Sscanf(price, "%f", &p)
			if err != nil {
				return nil, fmt.Errorf("failed to parse price: %w", err)
			}

			prices = append(prices, p)
		}
	}

	return prices, nil
}

func calculateSMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return nil
	}

	sma := make([]float64, len(prices)-period+1)

	for i := 0; i <= len(prices)-period; i++ {
		sum := 0.0
		for j := range period {
			sum += prices[i+j]
		}

		sma[i] = sum / float64(period)
	}

	return sma
}
