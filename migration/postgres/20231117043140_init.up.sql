CREATE TABLE profile (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    shedevrum_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TYPE profile_collector_type as ENUM('feed_top_day');

CREATE TABLE profile_collector (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    collector_type profile_collector_type NOT NULL UNIQUE,
    last_collected_at TIMESTAMP
);

INSERT INTO profile_collector (profile_collector_type) VALUES ('feed_top_day');

CREATE TABLE profile_metrics_collector (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    profile_handled_total INT NOT NULL,
    profile_handled_success INT NOT NULL,
    profile_handled_bad INT NOT NULL
);