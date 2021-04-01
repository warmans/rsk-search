CREATE TABLE "tscript_chunk_timeline"
(
    id          TEXT PRIMARY KEY,
    chunk_id    TEXT REFERENCES "tscript_chunk" (id) ON DELETE CASCADE,
    who         TEXT,
    what        TEXT,
    activity_at TIMESTAMP
);
