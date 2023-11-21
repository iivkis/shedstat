package domain

import "time"

type MetricsEntity struct {
	ProfileID     int       `db:"profile_id" json:"profile_id"`
	ShedevrumID   string    `db:"shedevrum_id" json:"shedevrum_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	Subscriptions uint64    `db:"subscriptions" json:"subscriptions"`
	Subscribers   uint64    `db:"subscribers" json:"subscribers"`
	Likes         uint64    `db:"likes" json:"likes"`
}

type MetricsChartEntity struct {
	Date          time.Time `db:"date" json:"date"`
	Subscriptions uint64    `db:"subscriptions" json:"subscriptions"`
	Subscribers   uint64    `db:"subscribers" json:"subscribers"`
	Likes         uint64    `db:"likes" json:"likes"`
}

type MetricsGetTopFilter string

const (
	METRICS_GET_TOP_FILTER_SUBSCRIPTIONS MetricsGetTopFilter = "subscriptions"
	METRICS_GET_TOP_FILTER_SUBSCRIBERS   MetricsGetTopFilter = "subscribers"
	METRICS_GET_TOP_FILTER_LIKES         MetricsGetTopFilter = "likes"
)
