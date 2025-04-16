CREATE TABLE transcript_tag
(
    episode_id TEXT    NOT NULL,
    tag_name   TEXT NOT NULL,
    tag_timestamp TEXT,
    PRIMARY KEY ("episode_id", "tag_name", "tag_timestamp")
);
