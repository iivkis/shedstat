package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func NewDBClickHouse(opts *clickhouse.Options) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(opts)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}
