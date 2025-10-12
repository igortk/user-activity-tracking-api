CREATE TABLE events (
    event_id   BIGSERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL CHECK (user_id > 0),
    event_action_timestamp  TIMESTAMPTZ NOT NULL,
    action     TEXT NOT NULL CHECK (action IN ('created', 'updated', 'deleted', 'viewed')),
    metadata   JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_user_time ON events(user_id, created_at);