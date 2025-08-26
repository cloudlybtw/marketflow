package storage

import (
	"context"
	"marketflow/internal/adapter/storage/cache"
	"marketflow/internal/adapter/storage/db"
)

type PriceRepository struct {
	Postgres *db.PostgresClient
	Redis    *cache.RedisClient
}

func NewPriceRepository(p *db.PostgresClient, r *cache.RedisClient) *PriceRepository {
	return &PriceRepository{
		Postgres: p,
		Redis:    r,
	}
}

func (p *PriceRepository) Health(ctx context.Context) error {
	err := p.Postgres.Ping()
	if err != nil {
		return err
	}

	err = p.Redis.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *PriceRepository) SaveLatest(ctx context.Context, exchange, pair string, price float64) {
	p.Redis.StorePriceUpdate(ctx, exchange, pair, price)
}

func (p *PriceRepository) GetLatest(ctx context.Context, symbol string) (float64, error) {
	return p.Redis.GetLatestPrice(ctx, "", symbol)
}

func (p *PriceRepository) GetLatestByExchange(ctx context.Context, exchange, symbol string) (float64, error) {
	return p.Redis.GetLatestPrice(ctx, exchange, symbol)
}
