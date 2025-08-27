package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"marketflow/internal/core/domain"
)

type PostgresClient struct {
	Db *sql.DB
}

func NewPostgresClient(db *sql.DB) *PostgresClient {
	return &PostgresClient{Db: db}
}

func (p *PostgresClient) Ping() error {
	return p.Db.Ping()
}

func (p *PostgresClient) SaveAggregatedData(data domain.AggregatedData) error {
	tx, err := postgresDB.Begin()
	if err != nil {
		return err
	}

	for key, exchange := range data.Exchanges {
		for pair, stat := range exchange.Coins {
			_, err = tx.Exec(`
				INSERT INTO prices (pair_name, exchange, timestamp, average_price, min_price, max_price) 
				VALUES ($1, $2, $3, $4, $5, $6);
				`, pair, key, stat.LastSeen, stat.AvgPrice, stat.MinPrice, stat.MaxPrice)
			if err != nil {
				slog.Error(err.Error())
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	slog.Info(fmt.Sprintf("Saved exchange stats for window %s to postgres",
		data.WindowEnd.Format("15:04:05")))
	return nil
}

func (p *PostgresClient) GetLowestByExchange(exchange, symbol string) (float64, error) {
	var minN float64
	err := p.Db.QueryRow(`
	SELECT min_price 
	FROM prices
	WHERE pair_name = $1
	AND exchange = $2
	ORDER BY timestamp DESC
	LIMIT 1;`, symbol, exchange).Scan(&minN)
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return minN, nil
}
