CREATE TABLE IF NOT EXISTS reminders (
    id TEXT PRIMARY KEY,
    application_id TEXT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    due_at TIMESTAMPTZ NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reminders_application_id
    ON reminders(application_id);

CREATE INDEX IF NOT EXISTS idx_reminders_due_at
    ON reminders(due_at);

CREATE INDEX IF NOT EXISTS idx_reminders_completed_due_at
    ON reminders(completed, due_at);