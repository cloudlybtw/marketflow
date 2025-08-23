package storage

type PriceRepository struct {
	Postgres *postgres.PostgresClient
	Redis    *redis.RedisClient
}

func NewPriceRepository(p *postgres.PostgresClient, r *redis.RedisClient) *PriceRepository {
	return &PriceRepository{
		Postgres: p,
		Redis:    r,
	}
}
