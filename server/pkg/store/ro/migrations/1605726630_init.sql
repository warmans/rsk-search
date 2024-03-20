CREATE TABLE "episode"
(
    "id"           TEXT PRIMARY KEY,
    "publication"  TEXT,
    "series"       INTEGER,
    "episode"      INTEGER,
    "release_date" TIMESTAMP,
    "metadata"     JSON,
    "contributors" JSON
);

CREATE TABLE "dialog"
(
    "id"              TEXT PRIMARY KEY,
    "episode_id"      TEXT REFERENCES episode (id),
    "pos"             INTEGER NOT NULL,
    "offset"          INTEGER NULL,
    "offset_inferred" BOOLEAN NOT NULL DEFAULT TRUE,
    "type"            TEXT    NOT NULL,
    "actor"           TEXT,
    "content"         TEXT,
    "metadata"        JSON,
    "notable"         BOOLEAN
);

CREATE INDEX dialog_pos ON dialog ("pos");

CREATE TABLE "changelog"
(
    "date"    DATE,
    "content" TEXT
);

CREATE INDEX changelog_date ON changelog ("date");

CREATE TABLE "song"
(
    "spotify_uri" TEXT,
    "artist"      TEXT,
    "title"       TEXT,
    "album"       TEXT,
    "episode_ids" JSONB
);

CREATE INDEX song_title ON song ("title");
CREATE INDEX song_artist ON song ("artist");
CREATE INDEX song_album ON song ("album");
