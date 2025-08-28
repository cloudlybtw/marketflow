package storage

import (
	"context"
	"fmt"
	"marketflow/internal/adapter/storage/cache"
	"marketflow/internal/adapter/storage/db"
	"marketflow/internal/core/domain"
	"time"
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

func (p *PriceRepository) SaveAggregated(data domain.AggregatedData) error {
	return p.Postgres.SaveAggregatedData(data)
}

func (p *PriceRepository) GetLatest(ctx context.Context, symbol string) (float64, string, error) {
	return p.Redis.GetLatestPrice(ctx, "", symbol)
}

func (p *PriceRepository) GetLatestByExchange(ctx context.Context, exchange, symbol string) (float64, string, error) {
	return p.Redis.GetLatestPrice(ctx, exchange, symbol)
}

func (p *PriceRepository) GetLowestByExchange(exchange, symbol string, duration time.Duration) (float64, string, int64, error) {

	if duration == 0 {
		return p.Postgres.GetLowestByExchange(exchange, symbol)
	}
	prices, err := p.Redis.GetPriceHistoryByExchange(symbol, exchange, time.Now().Add(-duration), time.Now())
	if err != nil || prices == nil {
		return -1, "", time.Now().Unix(), err
	}

	if prices == nil || len(prices) == 0 {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s on %s in last %v", symbol, exchange, duration)
	}
	var minN float64
	var tsN int64
	minN = 0

	for _, price := range prices {
		if minN == 0 || price.Price < minN {
			minN = price.Price
			tsN = price.Timestamp
		}
	}
	return minN, exchange, tsN, nil
}

func (p *PriceRepository) GetLowestGlobal(symbol string, duration time.Duration) (float64, string, int64, error) {
	if duration == 0 {

		return p.Postgres.GetLowestGlobal(symbol, duration)
	}

	exchanges := []string{"exchange1", "exchange2", "exchange3"}
	var minN float64
	var exN string
	var tsN int64
	for _, ex := range exchanges {
		prices, err := p.Redis.GetPriceHistoryByExchange(symbol, ex, time.Now().Add(-duration), time.Now())
		if err != nil || len(prices) == 0 {
			continue
		}
		for _, price := range prices {
			if minN == 0 || price.Price < minN {
				minN = price.Price
				exN = ex
				tsN = price.Timestamp
			}
		}
	}

	if minN == 0 {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s in last %v", symbol, duration)
	}
	return minN, exN, tsN, nil
}

func (p *PriceRepository) GetHighestByExchange(exchange, symbol string, duration time.Duration) (float64, string, int64, error) {
	if duration == 0 {
		return p.Postgres.GetHighestByExchange(exchange, symbol)
	}
	prices, err := p.Redis.GetPriceHistoryByExchange(symbol, exchange, time.Now().Add(-duration), time.Now())
	if err != nil || len(prices) == 0 {
		return -1, "", time.Now().Unix(), err
	}
	var maxN float64
	var exN string
	var tsN int64
	for _, price := range prices {
		if price.Price > maxN {
			maxN = price.Price
			exN = price.Exchange
			tsN = price.Timestamp
		}
	}
	return maxN, exN, tsN, nil
}

func (p *PriceRepository) GetHighestGlobal(symbol string, duration time.Duration) (float64, string, int64, error) {
	if duration == 0 {
		return p.Postgres.GetHighestGlobal(symbol, duration)
	}

	exchanges := []string{"exchange1", "exchange2", "exchange3"}
	var maxN float64
	var exN string
	var tsN int64
	for _, ex := range exchanges {
		prices, err := p.Redis.GetPriceHistoryByExchange(symbol, ex, time.Now().Add(-duration), time.Now())
		if err != nil || len(prices) == 0 {
			continue
		}
		for _, price := range prices {
			if price.Price > maxN {
				maxN = price.Price
				exN = price.Exchange
				tsN = price.Timestamp
			}
		}
	}

	if maxN == 0 {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s in last %v", symbol, duration)
	}

	return maxN, exN, tsN, nil
}

func (p *PriceRepository) GetAverageByExchange(exchange, symbol string, duration time.Duration) (float64, string, int64, error) {
	if duration == 0 {
		return p.Postgres.GetAverageByExchange(exchange, symbol)
	}
	prices, err := p.Redis.GetPriceHistoryByExchange(symbol, exchange, time.Now().Add(-duration), time.Now())
	if err != nil || len(prices) == 0 {
		return -1, "", time.Now().Unix(), err
	}
	var avgN float64
	var exN string
	var tsN int64
	for _, price := range prices {
		if price.Price > avgN {
			avgN = price.Price
		}
	}
	return avgN, exN, tsN, nil
}

func (p *PriceRepository) GetAverageGlobal(symbol string, duration time.Duration) (float64, string, int64, error) {
	if duration == 0 {
		return p.Postgres.GetAverageGlobal(symbol, duration)
	}

	exchanges := []string{"exchange1", "exchange2", "exchange3"}
	var avgN float64
	var exN string
	var tsN int64
	for _, ex := range exchanges {
		prices, err := p.Redis.GetPriceHistoryByExchange(symbol, ex, time.Now().Add(-duration), time.Now())
		if err != nil || len(prices) == 0 {
			continue
		}
		for _, price := range prices {
			if price.Price > avgN {
				avgN = price.Price
			}
		}
	}

	if avgN == 0 {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s in last %v", symbol, duration)
	}

	return avgN, exN, tsN, nil
}
