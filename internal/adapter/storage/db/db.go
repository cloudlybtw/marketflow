package db

import (
	"database/sql"
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
