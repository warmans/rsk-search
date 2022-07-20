CREATE TABLE "tscript_import"
(
    id           TEXT PRIMARY KEY,
    epid         TEXT,
    epname       TEXT NULL,
    mp3_uri      TEXT,
    log          JSONB,
    created_at   TIMESTAMP,
    completed_at TIMESTAMP
);

ALTER TABLE "tscript"
    ADD COLUMN name TEXT NULL;
