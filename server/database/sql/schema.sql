BEGIN;

CREATE TABLE IF NOT EXISTS post_briefs (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    slug        TEXT UNIQUE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category    TEXT,
    tags        TEXT[]
);

CREATE TABLE IF NOT EXISTS post_details (
    id          INTEGER REFERENCES post_briefs(id) ON DELETE CASCADE,
    content     TEXT
);

WITH ins AS (
    INSERT INTO post_briefs (title, slug)
    VALUES ('Home', 'home'), ('About', 'about')
    ON CONFLICT (slug) DO NOTHING
    RETURNING id
)
INSERT INTO post_details (id, content)
SELECT id, '' FROM ins;

CREATE TABLE IF NOT EXISTS comments (
    id          SERIAL PRIMARY KEY,
    post_id     INTEGER REFERENCES post_briefs(id) ON DELETE CASCADE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content     TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS local_users (
    username    TEXT PRIMARY KEY,
    password    TEXT NOT NULL,
    CHECK (username = 'admin')
);

COMMIT;