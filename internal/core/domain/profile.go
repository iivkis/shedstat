package domain

import (
	"time"
)

type ProfileEnity struct {
	ID            int       `db:"id" json:"id"`
	ShedevrumID   string    `db:"shedevrum_id" json:"shedevrum_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Link          string    `db:"link" json:"link"`
	AvatarURL     string    `db:"-" json:"avatar_url"`
	Name          string    `db:"-" json:"name"`
	Subscriptions uint64    `db:"-" json:"subscriptions"`
	Subscribers   uint64    `db:"-" json:"subscribers"`
	Likes         uint64    `db:"-" json:"likes"`
}
