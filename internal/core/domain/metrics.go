package domain

import "time"

type MetricEntity struct {
	ProfileID     int       `db:"profile_id" json:"profile_id"`
	ShedevrumID   string    `db:"shedevrum_id" json:"shedevrum_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Subscriptions uint64      `db:"subscriptions" json:"subscriptions"`
	Subscribers   uint64     `db:"subscribers" json:"subscribers"`
	Likes         uint64      `db:"likes" json:"likes"`
}
