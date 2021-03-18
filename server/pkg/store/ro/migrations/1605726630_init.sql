CREATE TABLE "episode"
(
    id           TEXT PRIMARY KEY,
    publication  TEXT,
    series       INTEGER,
    episode      INTEGER,
    release_date TIMESTAMP,
    metadata     JSON,
    tags         JSON
);

CREATE TABLE "dialog"
(
    id           TEXT PRIMARY KEY,
    episode_id   TEXT REFERENCES episode (id),
    pos          INTEGER NOT NULL,
    type         TEXT    NOT NULL,
    actor        TEXT,
    content      TEXT,
    metadata     JSON,
    content_tags JSON
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
    tscript_id   TEXT,
    raw          TEXT,
    start_second INT,
    end_second   INT
);
