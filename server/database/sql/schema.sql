BEGIN;

CREATE TABLE IF NOT EXISTS posts_brief (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category    TEXT,
    tags        TEXT[]
);

CREATE TABLE IF NOT EXISTS posts_detail (
    id          INTEGER REFERENCES posts_brief(id) ON DELETE CASCADE,
    content     TEXT
);

CREATE TABLE IF NOT EXISTS comments (
    id          SERIAL PRIMARY KEY,
    post_id     INTEGER REFERENCES posts_brief(id) ON DELETE CASCADE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content     TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS local_users (
    username    TEXT PRIMARY KEY,
    password    TEXT NOT NULL,
    CHECK (username = 'admin')
);

COMMIT;