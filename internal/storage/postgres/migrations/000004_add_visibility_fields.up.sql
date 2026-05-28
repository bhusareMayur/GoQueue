ALTER TABLE jobs
ADD COLUMN processing_started_at TIMESTAMP,
ADD COLUMN worker_id TEXT;