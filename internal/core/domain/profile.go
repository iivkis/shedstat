package domain

import (
	"database/sql"
	"time"
)

type ProfileEnity struct {
	ID            uint64    `db:"id" json:"-"`
	ShedevrumID   string    `db:"shedevrum_id" json:"shedevrum_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Link          string    `db:"link" json:"link"`
	Name          string    `db:"-" json:"name"`
	AvatarURL     string    `db:"-" json:"avatar_url"`
	Subscriptions uint64    `db:"-" json:"subscriptions"`
	Subscribers   uint64    `db:"-" json:"subscribers"`
	Likes         uint64    `db:"-" json:"likes"`
}

type ProfileMetricsCollectorEntity struct {
	ID                    uint64    `db:"id" json:"-"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	ProfileHandledTotal   uint64    `db:"profile_handled_total" json:"profile_handled_total"`
	ProfileHandledSuccess uint64    `db:"profile_handled_success" json:"profile_handled_success"`
	ProfileHandledBad     uint64    `db:"profile_handled_bad" json:"profile_handled_bad"`
}

type ProfileCollector_CollectorType string

const (
	PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY ProfileCollector_CollectorType = "feed_top_day"
)

type ProfileCollectorEntity struct {
	ID              uint8                          `db:"id" json:"-"`
	CollectorType   ProfileCollector_CollectorType `db:"collector_type" json:"collector_type"`
	LastCollectedAt sql.NullTime                   `db:"last_collected_at" json:"last_collected_at"`
}

type ProfileMetricsEntity struct {
	ProfileID     uint64    `db:"profile_id" json:"-"`
	ShedevrumID   string    `db:"shedevrum_id" json:"shedevrum_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Subscriptions uint64    `db:"subscriptions" json:"subscriptions"`
	Subscribers   uint64    `db:"subscribers" json:"subscribers"`
	Likes         uint64    `db:"likes" json:"likes"`
}

type ProfileMetrics_GetTopFilter string

const (
	PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIPTIONS ProfileMetrics_GetTopFilter = "subscriptions"
	PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIBERS   ProfileMetrics_GetTopFilter = "subscribers"
	PROFILE_METRICS_GET_TOP_FILTER_LIKES         ProfileMetrics_GetTopFilter = "likes"
)

func (f *ProfileMetrics_GetTopFilter) Scan(value string) {
	switch value {
	case string(PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIPTIONS):
		*f = PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIPTIONS
	case string(PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIBERS):
		*f = PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIBERS
	case string(PROFILE_METRICS_GET_TOP_FILTER_LIKES):
		*f = PROFILE_METRICS_GET_TOP_FILTER_LIKES
	default:
		*f = PROFILE_METRICS_GET_TOP_FILTER_SUBSCRIBERS
	}
}
