package core

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type DataCollectorService struct {
	wsClient   WebSocketClient
	repository OrderBookRepository
}

func NewDataCollectorService(wsClient WebSocketClient, repository OrderBookRepository) *DataCollectorService {
	return &DataCollectorService{
		wsClient:   wsClient,
		repository: repository,
	}
}

func (s *DataCollectorService) Run(ctx context.Context) error {
	// Setup graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Connect to WebSocket
	if err := s.wsClient.Connect(); err != nil {
		return err
	}
	defer s.wsClient.Close()

	log.Println("Successfully connected to WebSocket")

	// Start reading messages
	msgChan, errChan := s.wsClient.ReadMessages()

	for {
		select {
		case message := <-msgChan:
			log.Printf("Received data: %s", message)

			update := OrderBookUpdate{
				Data:      string(message),
				Timestamp: time.Now(),
			}

			if err := s.repository.Save(ctx, update); err != nil {
				log.Printf("Failed to save order book update: %v", err)
			}

		case err := <-errChan:
			return err

		case <-interrupt:
			log.Println("Termination signal received, shutting down...")

			return nil

		case <-ctx.Done():
			return nil
		}
	}
}
