CREATE TABLE dead_jobs (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    payload JSONB NOT NULL,
    retry_count INT NOT NULL,
    failed_at TIMESTAMP NOT NULL,
    last_error TEXT
);