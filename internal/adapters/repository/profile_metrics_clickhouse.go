package repository

import (
	"context"
	"fmt"
	"shedstat/internal/core/domain"
	"strings"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ProfileMetricsClickHouseRepository struct {
	db clickhouse.Conn
}

func NewProfileMetricsClickHouseRepository(db clickhouse.Conn) *ProfileMetricsClickHouseRepository {
	return &ProfileMetricsClickHouseRepository{
		db: db,
	}
}

func (r *ProfileMetricsClickHouseRepository) Create(ctx context.Context, metrics []*domain.ProfileMetricsEntity) error {
	q := `
	INSERT INTO
		profile_metrics 
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

func (r *ProfileMetricsClickHouseRepository) GetByShedevrumID(shedevrumID string) ([]*domain.ProfileMetricsEntity, error) {
	q := `
		SELECT 
			toDate(created_at),
			MIN(subscriptions),
			MIN(subscribers),
			MIN(likes)
		FROM profile_metrics 
		WHERE shedevrum_id = $1
		GROUP BY toDate(created_at)
		ORDER BY toDate(created_at)
    `
	rows, err := r.db.Query(context.Background(), q, shedevrumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metrics := make([]*domain.ProfileMetricsEntity, 0)
	for rows.Next() {
		m := &domain.ProfileMetricsEntity{
			ShedevrumID: shedevrumID,
		}
		if err := rows.Scan(&m.CreatedAt, &m.Subscriptions, &m.Subscribers, &m.Likes); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (r *ProfileMetricsClickHouseRepository) GetTop(ctx context.Context, filter domain.ProfileMetrics_GetTopFilter, amount int) ([]*domain.ProfileMetricsEntity, error) {
	q := `SELECT DISTINCT ON(profile_id) * FROM profile_metrics ORDER BY toDate(created_at) DESC, ` + string(filter) + ` DESC LIMIT $1`
	rows, err := r.db.Query(ctx, q, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metrics := make([]*domain.ProfileMetricsEntity, 0, amount)
	for rows.Next() {
		m := &domain.ProfileMetricsEntity{}
		if err := rows.Scan(
			&m.ProfileID,
			&m.ShedevrumID,
			&m.Subscriptions,
			&m.Subscribers,
			&m.Likes,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}
