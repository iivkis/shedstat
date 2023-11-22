package repository

import (
	"context"
	"database/sql"
	"errors"
	"shedstat/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type ProfilePostgresRepository struct {
	db *sqlx.DB
}

func NewProfilePostgresRepository(db *sqlx.DB) *ProfilePostgresRepository {
	return &ProfilePostgresRepository{
		db: db,
	}
}

func (r *ProfilePostgresRepository) Create(ctx context.Context, shedevrumID string) error {
	q := `
		INSERT INTO
		profile (
			shedevrum_id
		)
		VALUES ($1)
	`
	_, err := r.db.ExecContext(ctx, q, shedevrumID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfilePostgresRepository) GetByID(ctx context.Context, id int) (*domain.ProfileEnity, error) {
	q := `SELECT * FROM profile WHERE id = $1`
	var p domain.ProfileEnity
	err := r.db.GetContext(ctx, &p, q, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProfilePostgresRepository) GetByShedevrumID(ctx context.Context, shedevrumID string) (*domain.ProfileEnity, error) {
	q := `SELECT * FROM profile WHERE shedevrum_id = $1`
	var p domain.ProfileEnity
	err := r.db.GetContext(ctx, &p, q, shedevrumID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProfilePostgresRepository) GetList(ctx context.Context, startFromID int, amount int) ([]*domain.ProfileEnity, error) {
	q := `SELECT * FROM profile WHERE id > $1 LIMIT $2`
	var profiles []*domain.ProfileEnity
	err := r.db.SelectContext(ctx, &profiles, q, startFromID, amount)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *ProfilePostgresRepository) ExistsByShedevrumID(ctx context.Context, shedevrumID string) (bool, error) {
	q := `SELECT EXISTS(SELECT id FROM profile WHERE shedevrum_id = $1)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, q, shedevrumID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
