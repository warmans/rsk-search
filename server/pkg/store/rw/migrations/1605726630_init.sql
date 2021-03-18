-- once a chunk has been completed it becomes a contribution
CREATE TABLE "tscript_contribution"
(
    id               TEXT PRIMARY KEY,
    author_id        TEXT,
    tscript_chunk_id TEXT,
    transcription    TEXT
);

CREATE TABLE "tscript_chunk_activity"
(
    chunk_id       TEXT PRIMARY KEY,
    last_requested TEXT NULL
);

CREATE TABLE "author"
(
    id   TEXT PRIMARY KEY,
    name TEXT
);
