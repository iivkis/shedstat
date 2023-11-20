package repository

import (
	"context"
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

func (r *ProfilePostgresRepository) Create(ctx context.Context, profile *domain.ProfileEnity) error {
	q := `
		INSERT INTO
		profile (
			shedevrum_id,
			link
		)
		VALUES ($1, $2)
	`
	var p domain.ProfileEnity
	err := r.db.GetContext(ctx, &p, q, profile.ShedevrumID, profile.Link)
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
	q := `SELECT EXISTS(*) FROM profile WHERE shedevrum_id = $1`
	var exists bool
	err := r.db.GetContext(ctx, &exists, q, shedevrumID)
	if err != nil {
		return false, err
	}
	return exists, nil
}
