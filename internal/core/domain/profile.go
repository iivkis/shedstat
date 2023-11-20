package domain

import (
	"time"
)

type ProfileEnity struct {
	ID          int       `db:"id" json:"id"`
	ShedevrumID string    `db:"shedevrum_id" json:"shedevrum_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Link        string    `db:"link" json:"link"`
}
