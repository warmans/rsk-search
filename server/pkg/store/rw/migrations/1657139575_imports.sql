CREATE TABLE "tscript_import"
(
    id      TEXT PRIMARY KEY,
    epid    TEXT,
    -- should be something like create_wav -> auto_transcribe -> create_chunks -> split_mp3 -> publish_chunks
    stage   TEXT,
    mp3_uri TEXT,
    log     JSONB
);
