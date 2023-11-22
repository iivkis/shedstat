CREATE TABLE profile_metrics (
    profile_id UInt64,
    shedevrum_id VARCHAR(255),
    subscriptions UInt64,
    subscribers UInt64,
    likes UInt64,
    created_at DateTime
) ENGINE = MergeTree PARTITION BY (created_at)
    TTL created_at + INTERVAL '3 month'
    ORDER BY (created_at)