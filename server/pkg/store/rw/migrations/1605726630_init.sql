CREATE TABLE "author"
(
    id         TEXT PRIMARY KEY,
    name       TEXT UNIQUE NOT NULL,
    identity   JSON,
    created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
    banned     BOOLEAN     NOT NULL DEFAULT false,
    approver   BOOLEAN     NOT NULL DEFAULT false
);

-- tscript is an incomplete transcription
CREATE TABLE "tscript"
(
    id          TEXT PRIMARY KEY,
    publication TEXT,
    series      INTEGER,
    episode     INTEGER
);

-- scripts are a series of audio chunks that have been auto-transcribed
CREATE TABLE "tscript_chunk"
(
    id           TEXT PRIMARY KEY,
    tscript_id   TEXT REFERENCES "tscript" (id) ON DELETE CASCADE,
    raw          TEXT,
    start_second INT,
    end_second   INT
);

-- once a chunk has been completed it becomes a contribution
-- do not delete them with the other data to ensure they are not lost
-- under most common errors
CREATE TABLE "tscript_contribution"
(
    id               TEXT PRIMARY KEY,
    author_id        TEXT,
    tscript_chunk_id TEXT,
    transcription    TEXT,
    state            TEXT,
    created_at       TIMESTAMP
);

CREATE TABLE "tscript_chunk_activity"
(
    tscript_chunk_id TEXT PRIMARY KEY REFERENCES tscript_chunk (id) ON DELETE CASCADE,
    accessed_at      TIMESTAMP NULL,
    submitted_at     TIMESTAMP NULL,
    approved_at      TIMESTAMP NULL,
    rejected_at      TIMESTAMP NULL
);
