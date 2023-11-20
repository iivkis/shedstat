CREATE TABLE profile (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    shedevrum_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    link VARCHAR(255)
)

CREATE TABLE metrics_shedule (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    profile_handled_total INT NOT NULL,
    profile_handled_success INT NOT NULL,
    profile_handled_bad INT NOT NULL
)