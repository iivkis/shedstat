package repository

import (
	"context"
	"fmt"
	"shedstat/internal/core/domain"
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

func (r *MetricsClickHouseRepository) Create(ctx context.Context, metrics []*domain.MetricsEntity) error {
	q := `
	INSERT INTO
		metrics 
	VALUES
	`
	values := make([]string, 0, len(metrics))
	for _, m := range metrics {
		values = append(values,
			fmt.Sprintf("(%d, '%s', %d, %d, %d, NOW())", m.ProfileID, m.ShedevrumID, m.Subscriptions, m.Subscribers, m.Likes),
		)
	}
	err := r.db.Exec(ctx, q+strings.Join(values, ", "))
	if err != nil {
		return err
	}
	return nil
}
