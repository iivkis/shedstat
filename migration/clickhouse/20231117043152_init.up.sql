CREATE TABLE metrics (
    profile_id INT,
    shedevrum_id VARCHAR(255),
    subscriptions UINT64,
    subscribers UINT64,
    likes UINT64,
    created_at DateTime
) ENGINE = MergeTree PARTITION BY (created_at)
    TTL created_at + INTERVAL '3 month'
    ORDER BY (created_at)