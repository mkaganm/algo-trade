package ports

import (
	"github.com/mkaganm/algo-trade/collector/internal/core"
)

type OrderBookRepository interface {
	core.OrderBookRepository
}
