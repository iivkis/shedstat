CREATE TABLE profile (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    shedevrum_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    link VARCHAR(255)
)

CREATE collector_type as ENUM(
    'feed_top_day'
);

CREATE TABLE profile_collector (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    collector_type collector_type NOT NULL UNIQUE,
    last_collected_at TIMESTAMP
);

INSERT INTO profile_collector (collector_type) VALUES ('feed_top_day');

CREATE TABLE metrics_collector (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    profile_handled_total INT NOT NULL,
    profile_handled_success INT NOT NULL,
    profile_handled_bad INT NOT NULL
);