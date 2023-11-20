package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBPostgresSource struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	SSLMode  string
}

func NewDBPostgres(s DBPostgresSource) (*sqlx.DB, error) {
	source := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		s.User, s.Password, s.Host, s.Port, s.Name)
	db, err := sqlx.Connect("pgx", source)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
