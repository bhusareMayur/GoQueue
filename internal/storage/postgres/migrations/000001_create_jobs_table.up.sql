CREATE TABLE jobs (
    id UUID PRIMARY KEY,

    type TEXT NOT NULL,

    payload JSONB NOT NULL,

    status TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);