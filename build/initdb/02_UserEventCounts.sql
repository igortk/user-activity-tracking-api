CREATE TABLE user_event_counts (
    user_id BIGINT NOT NULL,
    period_start TIMESTAMP NOT NULL,
    event_count INT NOT NULL DEFAULT 0,
    period_end TIMESTAMP NOT NULL,

    PRIMARY KEY (user_id, period_start)
);