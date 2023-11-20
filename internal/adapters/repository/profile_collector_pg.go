package repository

import (
	"database/sql"
	"shedstat/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type ProfileCollectorPostgresRepository struct {
	db *sqlx.DB
}

func NewProfileCollectorPostgresRepository(db *sqlx.DB) *ProfileCollectorPostgresRepository {
	return &ProfileCollectorPostgresRepository{
		db: db,
	}
}

func (r *ProfileCollectorPostgresRepository) GetLastCollectedAt(collectorType domain.ProfileCollector_CollectorType) (sql.NullTime, error) {
	q := `SELECT last_collected_at FROM profile_collector WHERE collector_type = $1`
	var lastCollectedAt sql.NullTime
	if err := r.db.Get(&lastCollectedAt, q, collectorType); err != nil {
		return sql.NullTime{}, err
	}
	return lastCollectedAt, nil
}

func (r *ProfileCollectorPostgresRepository) UpdateLastCollectedAt(collectorType domain.ProfileCollector_CollectorType) error {
	q := `UPDATE profile_collector SET last_collected_at = NOW() WHERE collector_type = $1`
	if _, err := r.db.Exec(q, collectorType); err != nil {
		return err
	}
	return nil
}
