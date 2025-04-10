package ports

import "context"

type HealthService interface {
	CheckHealth(ctx context.Context) error
}
