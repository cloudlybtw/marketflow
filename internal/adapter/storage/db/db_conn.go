package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"marketflow/internal/adapter/config"

	_ "github.com/lib/pq"
)

var postgresDB *sql.DB

func OpenDB(config *config.DB) (*sql.DB, error) {

	// postgres://username:password@host:port/dbname
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Name)

	if postgresDB != nil {
		return postgresDB, nil
	}

	slog.Info("Trying to connect to PostgreSQL database...")

	maxRetries := 6
	retryInterval := 5 * time.Second
	var db *sql.DB
	var err error

	// Try connecting up to maxRetries times
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			// Create a context with a timeout for the Ping operation
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Try to ping the database to check if it's available
			err = db.PingContext(ctx)
			if err == nil {
				// If the ping is successful, use the connection
				postgresDB = db
				slog.Info("PostgreSQL database connection established")
				return postgresDB, nil
			}
		}

		// If any error occurs, log it and retry after a delay
		slog.Error("Failed to connect to PostgreSQL, retrying in 5 seconds... (attempt ", strconv.Itoa(i+1), "/", strconv.Itoa(maxRetries), ")")
		time.Sleep(retryInterval)
	}

	slog.Error("Failed to connect to PostgreSQL after", "attempts", maxRetries, "interval", retryInterval, "error", err)
	return nil, err
}
