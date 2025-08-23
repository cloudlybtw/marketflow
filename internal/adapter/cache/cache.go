package cache

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"marketflow/internal/adapter/config"

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

func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}
