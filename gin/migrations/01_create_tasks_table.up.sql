CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    done BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL
);