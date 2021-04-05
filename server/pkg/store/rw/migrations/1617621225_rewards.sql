CREATE TABLE "author_reward"
(
    id         TEXT PRIMARY KEY,
    author_id  TEXT REFERENCES "author" (id) ON DELETE CASCADE,
    threshold  INTEGER,
    created_at TIMESTAMP,
    confirmed  BOOLEAN DEFAULT FALSE,
    error      TEXT NULL

);

CREATE UNIQUE INDEX author_threshold ON "author_reward" (author_id, threshold)
