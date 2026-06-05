CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    priority TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP
);

-- Index to optimize the background worker's polling queries
CREATE INDEX idx_outbox_events_status_created_at ON outbox_events(status, created_at);