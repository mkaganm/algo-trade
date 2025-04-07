package helpers

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func Retry(maxRetries int, delay time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	var err error

	var result interface{}

	for i := range maxRetries {
		result, err = fn()
		if err == nil {
			return result, nil
		}

		log.Printf("Attempt %d/%d failed: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	return nil, err
}

func RetryWebSocket(maxRetries int, delay time.Duration, url string) (*websocket.Conn, error) {
	var conn *websocket.Conn

	var err error

	for i := range maxRetries {
		conn, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err == nil {
			return conn, nil
		}

		log.Printf("Connection attempt %d/%d failed: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	return nil, err
}
