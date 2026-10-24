CREATE TABLE IF NOT EXISTS documents (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL CHECK (btrim(title) <> ''),
    content TEXT NOT NULL CHECK (btrim(content) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
