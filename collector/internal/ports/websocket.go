package ports

import "github.com/mkaganm/algo-trade/collector/internal/core"

type WebSocketClient interface {
	core.WebSocketClient
}
