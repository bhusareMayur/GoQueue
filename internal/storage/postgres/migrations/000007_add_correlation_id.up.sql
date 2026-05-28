ALTER TABLE jobs
ADD COLUMN correlation_id TEXT;

ALTER TABLE dead_jobs
ADD COLUMN correlation_id TEXT;