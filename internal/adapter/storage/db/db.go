package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"marketflow/internal/core/domain"
	"time"
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

func (p *PostgresClient) GetLowestByExchange(exchange, symbol string) (float64, string, int64, error) {
	var minN float64
	var exN string
	var tsN int64
	err := p.Db.QueryRow(`
	SELECT min_price, exchange, timestamp, 
	FROM prices
	WHERE pair_name = $1
	AND exchange = $2
	ORDER BY timestamp DESC
	LIMIT 1;`, symbol, exchange).Scan(&minN, &exN, &tsN)
	if err != nil {
		slog.Error(err.Error())
		return -1, "", time.Now().Unix(), err
	}

	return minN, exN, tsN, nil
}

func (p *PostgresClient) GetLowestGlobal(symbol string, duration time.Duration) (float64, string, int64, error) {
	var minN sql.NullFloat64
	var exN string
	var tsN int64
	query := `
        SELECT min_price, exchange, timestamp, 
		FROM prices
		WHERE pair_name = $1
    `
	args := []interface{}{symbol}

	if duration > 0 {
		query += " AND timestamp >= NOW() - $2::interval"
		args = append(args, fmt.Sprintf("%f seconds", duration.Seconds()))
	}

	err := p.Db.QueryRow(query+`
		ORDER BY min_price ASC
		LIMIT 1`, args...).Scan(&minN, &exN, &tsN)
	if err != nil {
		slog.Error("query error", slog.String("err", err.Error()))
		return -1, "", time.Now().Unix(), err
	}

	if !minN.Valid {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s in last %v", symbol, duration)
	}

	return minN.Float64, exN, tsN, nil
}

func (p *PostgresClient) GetHighestByExchange(exchange, symbol string) (float64, string, int64, error) {
	var maxN float64
	var exN string
	var tsN int64
	err := p.Db.QueryRow(`
        SELECT max_price, exchange, timestamp
        FROM prices
        WHERE pair_name = $1
          AND exchange = $2
        ORDER BY timestamp DESC
        LIMIT 1;`, symbol, exchange).Scan(&maxN, &exN, &tsN)

	if err != nil {
		slog.Error(err.Error())
		return -1, "", time.Now().Unix(), err
	}

	return maxN, exN, tsN, nil
}

func (p *PostgresClient) GetHighestGlobal(symbol string, duration time.Duration) (float64, string, int64, error) {
	var maxN sql.NullFloat64
	var exN string
	var tsN int64
	query := `
        SELECT max_price, exchange, timestamp, 
		FROM prices
		WHERE pair_name = $1
    `
	args := []interface{}{symbol}

	if duration > 0 {
		query += " AND timestamp >= NOW() - $2::interval"
		args = append(args, fmt.Sprintf("%f seconds", duration.Seconds()))
	}

	err := p.Db.QueryRow(query+`
		ORDER BY max_price DESC
		LIMIT 1`, args...).Scan(&maxN, &exN, &tsN)
	if err != nil {
		slog.Error("query error", slog.String("err", err.Error()))
		return -1, "", time.Now().Unix(), err
	}

	if !maxN.Valid {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s in last %v", symbol, duration)
	}

	return maxN.Float64, exN, tsN, nil
}

func (p *PostgresClient) GetAverageByExchange(exchange, symbol string) (float64, string, int64, error) {
	var avgN float64
	var exN string
	var tsN int64
	err := p.Db.QueryRow(`
        SELECT average_price, exchange, timestamp
        FROM prices
        WHERE pair_name = $1
          AND exchange = $2
        ORDER BY timestamp DESC
        LIMIT 1;`, symbol, exchange).Scan(&avgN, &exN, &tsN)

	if err != nil {
		slog.Error(err.Error())
		return -1, "", time.Now().Unix(), err
	}

	return avgN, exN, tsN, nil
}

func (p *PostgresClient) GetAverageGlobal(symbol string, duration time.Duration) (float64, string, int64, error) {
	var avgN sql.NullFloat64
	var exN string
	var tsN int64
	query := `
        SELECT average_price, exchange, timestamp
        FROM prices
        WHERE pair_name = $1
    `
	args := []interface{}{symbol}

	if duration > 0 {
		query += " AND timestamp >= NOW() - $2::interval"
		args = append(args, fmt.Sprintf("%f seconds", duration.Seconds()))
	}

	err := p.Db.QueryRow(query+`
        ORDER BY average_price DESC
        LIMIT 1`, args...).Scan(&avgN, &exN, &tsN)
	if err != nil {
		slog.Error("query error", slog.String("err", err.Error()))
		return -1, "", time.Now().Unix(), err
	}

	if !avgN.Valid {
		return -1, "", time.Now().Unix(), fmt.Errorf("no prices found for %s in last %v", symbol, duration)
	}

	return avgN.Float64, exN, tsN, nil
}
