package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	timeout = 2 * time.Second
)

type Response struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

func CheckHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		// Check MongoDB connection
		if err := client.Ping(ctx, nil); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(Response{
				Status:  "unhealthy",
				Details: "Cannot connect to MongoDB",
			})

			return
		}

		// If all checks pass
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(Response{
			Status:  "healthy",
			Details: "All systems operational",
		})
	}
}
