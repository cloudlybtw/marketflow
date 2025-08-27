package port

import (
	"context"
	"marketflow/internal/core/domain"
	"time"
)

type PriceRepository interface {
	SaveAggregated(data domain.AggregatedData) error
	SaveLatest(ctx context.Context, exchange, symbol string, price float64)
	GetLatest(ctx context.Context, symbol string) (float64, error)
	GetLatestByExchange(ctx context.Context, exchange, symbol string) (float64, error)
	GetLowestByExchange(exchange, symbol string, duration time.Duration) (float64, error)
	Health(ctx context.Context) error
}

type PriceService interface {
	SavePrice(ctx context.Context, in <-chan domain.Message)
	GetLatest(ctx context.Context, exchange string, pair string) (*domain.PriceUpdate, error)
	// GetHighest()
	GetLowest(exchange, pair string, duration time.Duration) (*domain.PriceUpdate, error)
	// GetAverage()
	CheckHealth(ctx context.Context) error
}
