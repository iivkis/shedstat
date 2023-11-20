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

func (r *MetricsClickHouseRepository) GetByShedevrumID(shedevrumID string) ([]*domain.MetricsProfileEntity, error) {
	q := `
		SELECT 
			toDate(created_at) as date,
			MIN(subscriptions) as subscriptions,
			MIN(subscribers) as subscribers,
			MIN(likes) as likes
		FROM metrics 
		WHERE shedevrum_id = '%s'
		GROUP BY toDate(created_at)
		ORDER BY toDate(created_at)
    `
	rows, err := r.db.Query(context.Background(), q, shedevrumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metrics := make([]*domain.MetricsProfileEntity, 0)
	for rows.Next() {
		m := &domain.MetricsProfileEntity{}
		err := rows.Scan(&m.Date, &m.Subscriptions, &m.Subscribers, &m.Likes)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
