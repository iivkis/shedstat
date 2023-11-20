package domain

import "time"

type MetricSheduleEntity struct {
	ID                    int       `db:"id" json:"id"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	ProfileHandledTotal   uint64    `db:"profile_handled_total" json:"profile_handled_total"`
	ProfileHandledSuccess uint64    `db:"profile_handled_success" json:"profile_handled_success"`
	ProfileHandledBad     uint64    `db:"profile_handled_bad" json:"profile_handled_bad"`
}
