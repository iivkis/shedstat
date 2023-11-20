package repository

import (
	"context"
	"shedstat/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type MetricsShedulePostgresRepository struct {
	db *sqlx.DB
}

func NewMetricsShedulePostgresRepository(db *sqlx.DB) *MetricsShedulePostgresRepository {
	return &MetricsShedulePostgresRepository{
		db: db,
	}
}

func (r *MetricsShedulePostgresRepository) Create(ctx context.Context, entity *domain.MetricSheduleEntity) error {
	q := `
		INSERT INTO metrics_schedule (
            profile_handled_total,
            profile_handled_success,
            profile_handled_bad
		) VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, q, entity.ProfileHandledTotal, entity.ProfileHandledSuccess, entity.ProfileHandledBad)
	if err != nil {
		return err
	}
	return nil
}

func (r *MetricsShedulePostgresRepository) GetLast(ctx context.Context) (*domain.MetricSheduleEntity, error) {
	q := `
		SELECT * FROM metrics_schedule ORDER BY id DESC LIMIT 1
	`
	var entity domain.MetricSheduleEntity
	if err := r.db.GetContext(ctx, &entity, q); err != nil {
		return nil, err
	}
	return &entity, nil
}
