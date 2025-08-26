package port

import (
	"context"
	"marketflow/internal/core/domain"
)

type PriceRepository interface {
	// SaveAggregated(data domain.AggregatedData) error
	SaveLatest(ctx context.Context, exchange, symbol string, price float64)
	Health(ctx context.Context) error
}

type PriceService interface {
	SavePrice(ctx context.Context, in <-chan domain.Message)
	// GetLatest(ctx context.Context, exchange string, pair string)
	// GetHighest()
	// GetLowest()
	// GetAverage()
	CheckHealth(ctx context.Context) error
}
