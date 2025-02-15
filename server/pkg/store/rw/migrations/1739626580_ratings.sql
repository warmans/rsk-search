ALTER TABLE author
    ADD COLUMN IF NOT EXISTS placeholder BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE transcript_rating_score
(
    author_id  TEXT    NOT NULL REFERENCES author (id) ON DELETE CASCADE,
    episode_id TEXT    NOT NULL,
    score      NUMERIC NOT NULL DEFAULT 0,
    -- since these get synced to flat files there needs to be a way to indicate a delete operation
    delete     BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY ("author_id", episode_id)
);
