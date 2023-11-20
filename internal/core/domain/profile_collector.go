package domain

import (
	"database/sql"
)

type ProfileCollector_CollectorType = string

const (
	PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY = "feed_top_day"
)

type ProfileCollectorEntity struct {
	ID              int                            `db:"id" json:"id"`
	CollectorType   ProfileCollector_CollectorType `db:"collector_type" json:"collector_type"`
	LastCollectedAt sql.NullTime                   `db:"last_collected_at" json:"last_collected_at"`
}
