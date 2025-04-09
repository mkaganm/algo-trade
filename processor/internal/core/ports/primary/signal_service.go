package primary

import (
	"github.com/mkaganm/algo-trade/processor/internal/core/domain"
)

// SignalService is the primary port (interface) for signal processing
type SignalService interface {
	GenerateSignal() (*domain.TradeSignal, error)
	CalculateSMAs(prices []float64) (shortSMA, longSMA []float64, err error)
}
