CREATE TABLE IF NOT EXISTS application_status_history (
    id BIGSERIAL PRIMARY KEY,
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    from_status TEXT NOT NULL,
    to_status TEXT NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_application_status_history_application_id
    ON application_status_history(application_id);

CREATE INDEX IF NOT EXISTS idx_application_status_history_changed_at
    ON application_status_history(changed_at);

CREATE TABLE IF NOT EXISTS activity_events (
    id BIGSERIAL PRIMARY KEY,
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    description TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_activity_events_application_id
    ON activity_events(application_id);

CREATE INDEX IF NOT EXISTS idx_activity_events_occurred_at
    ON activity_events(occurred_at);