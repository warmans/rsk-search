CREATE TABLE "transcript_change"
(
    id            TEXT PRIMARY KEY,
    author_id     TEXT,
    epid          TEXT,
    summary       TEXT,
    transcription TEXT,
    state         TEXT,
    created_at    TIMESTAMP
);

DROP TABLE tscript_chunk_timeline;
