package repository

import (
	"context"
	"shedstat/internal/core/domain"
	"strconv"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type MetricsClickHouseRepository struct {
	db clickhouse.Conn
}

func NewMetricsClickHouseRepository(db clickhouse.Conn) *MetricsClickHouseRepository {
	return &MetricsClickHouseRepository{
		db: db,
	}
}

func (r *MetricsClickHouseRepository) Create(ctx context.Context, metrics []*domain.MetricEntity) error {
	q := `
	INSERT INTO
		metrics 
	VALUES
	`
	values := make([]string, 0, len(metrics))
	for _, m := range metrics {
		values = append(values, "("+
			m.ShedevrumID+","+
			strconv.FormatUint(m.Subscriptions, 10)+","+
			strconv.FormatUint(m.Subscribers, 10)+","+
			strconv.FormatUint(m.Likes, 10)+","+
			"NOW()"+
			")",
		)
	}
	err := r.db.Exec(ctx, q+strings.Join(values, ", "))
	if err != nil {
		return err
	}
	return nil
}
