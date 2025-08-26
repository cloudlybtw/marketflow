package cache

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"marketflow/internal/adapter/config"
	"marketflow/internal/core/domain"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisClient(r *config.Redis) (*RedisClient, error) {
	db, _ := strconv.Atoi(r.Db)

	client := redis.NewClient(&redis.Options{
		Addr:     r.Host + ":" + r.Port,
		Password: r.Password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &RedisClient{client: client}, nil
}

func (r *RedisClient) StorePriceUpdate(ctx context.Context, exchange, pair string, price float64) {
	latestKeyByExchange := "L:" + exchange + ":" + pair
	if err := r.client.Set(ctx, latestKeyByExchange, price, 5*time.Minute).Err(); err != nil {
		r.logger.Error("failed to set key", "key", latestKeyByExchange, "error", err)
		return
	}
	latestKey := "L:" + pair
	if err := r.client.Set(ctx, latestKey, price, 5*time.Minute).Err(); err != nil {
		r.logger.Error("failed to set key", "key", latestKey, "error", err)
		return
	}
	historyKey := "H:" + pair
	r.AddToSortedSet(ctx, historyKey, float64(time.Now().Unix()), price)

	historyExchangeKey := "H:" + exchange + ":" + pair
	r.AddToSortedSet(ctx, historyExchangeKey, float64(time.Now().Unix()), price)
}

func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *RedisClient) AddToSortedSet(ctx context.Context, key string, score float64, member interface{}) error {
	return r.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

func (r *RedisClient) GetPriceHistory(pair string, from, to time.Time) ([]domain.PriceUpdate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := "H:" + pair
	results, err := r.client.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(from.Unix(), 10),
		Max: strconv.FormatInt(to.Unix(), 10),
	}).Result()

	var prices []domain.PriceUpdate
	for _, z := range results {
		if price, ok := z.Member.(string); ok {
			if p, err := strconv.ParseFloat(price, 64); err == nil {
				prices = append(prices, domain.PriceUpdate{
					Symbol:    pair,
					Timestamp: int64(z.Score),
					Price:     p,
				})
			}
		}
	}
	return prices, err
}

func (r *RedisClient) GetPriceHistoryByExchange(pair, exchange string, from, to time.Time) ([]domain.PriceUpdate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	key := "H:" + exchange + ":" + pair
	results, err := r.client.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(from.Unix(), 10),
		Max: strconv.FormatInt(to.Unix(), 10),
	}).Result()

	var prices []domain.PriceUpdate
	for _, z := range results {
		if price, ok := z.Member.(string); ok {
			if p, err := strconv.ParseFloat(price, 64); err == nil {
				prices = append(prices, domain.PriceUpdate{
					Symbol:    pair,
					Timestamp: int64(z.Score),
					Price:     p,
				})
			}
		}
	}
	return prices, err
}
