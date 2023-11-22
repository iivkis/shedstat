package repository

import (
	"context"
	"shedstat/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type MetricsCollectorPostgresRepository struct {
	db *sqlx.DB
}

func NewMetricsCollectorPostgresRepository(db *sqlx.DB) *MetricsCollectorPostgresRepository {
	return &MetricsCollectorPostgresRepository{
		db: db,
	}
}

func (r *MetricsCollectorPostgresRepository) Create(ctx context.Context, entity *domain.MetricsCollectorEntity) error {
	q := `
		INSERT INTO profile_metrics_collector (
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

func (r *MetricsCollectorPostgresRepository) GetLast(ctx context.Context) (*domain.MetricsCollectorEntity, error) {
	q := `
		SELECT * FROM profile_metrics_collector ORDER BY id DESC LIMIT 1
	`
	var entity domain.MetricsCollectorEntity
	if err := r.db.GetContext(ctx, &entity, q); err != nil {
		return nil, err
	}
	return &entity, nil
}
