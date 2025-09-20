BEGIN;

CREATE TABLE IF NOT EXISTS post_briefs (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    slug        TEXT UNIQUE,
    category    TEXT,
    tags        TEXT[],
    status      TEXT NOT NULL,
    pinned      BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (status IN ('draft', 'published', 'deleted', 'hidden'))
);

CREATE TABLE IF NOT EXISTS post_details (
    id          INTEGER REFERENCES post_briefs(id) ON DELETE CASCADE,
    content     TEXT,
    CONSTRAINT unique_post_id UNIQUE (id)
);

WITH ins AS (
    INSERT INTO post_briefs (title, slug, status)
    VALUES ('Home', 'home', 'hidden'), ('About', 'about', 'hidden')
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