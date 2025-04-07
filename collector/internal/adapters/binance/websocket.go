package binance

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/mkaganm/algo-trade/collector/internal/helpers"
)

type WebSocket struct {
	url        string
	conn       *websocket.Conn
	maxRetries int
	retryDelay time.Duration
}

func NewBinanceWebSocket(url string, maxRetries int, retryDelay time.Duration) *WebSocket {
	return &WebSocket{
		url:        url,
		conn:       nil,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

func (b *WebSocket) Connect() error {
	conn, err := helpers.RetryWebSocket(b.maxRetries, b.retryDelay, b.url)
	if err != nil {
		return err
	}

	b.conn = conn

	return nil
}

func (b *WebSocket) ReadMessages() (<-chan []byte, <-chan error) {
	msgChan := make(chan []byte)
	errChan := make(chan error)

	go b.readMessagesGoroutine(msgChan, errChan)

	return msgChan, errChan
}

func (b *WebSocket) readMessagesGoroutine(msgChan chan<- []byte, errChan chan<- error) {
	defer close(msgChan)
	defer close(errChan)
	defer helpers.RecoverRoutine(errChan)

	for {
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				errChan <- err
			}

			return
		}

		msgChan <- message
	}
}

func (b *WebSocket) Close() error {
	if b.conn != nil {
		return b.conn.Close()
	}

	return nil
}
