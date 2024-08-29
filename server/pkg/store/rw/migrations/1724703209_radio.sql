CREATE TABLE "radio_episode"
(
    episode_id  TEXT PRIMARY KEY,
    publication TEXT
);

CREATE TABLE radio_exclusion
(
    author_id  TEXT REFERENCES "author" (id) ON DELETE CASCADE,
    episode_id TEXT REFERENCES "radio_episode" (episode_id) ON DELETE CASCADE,
    PRIMARY KEY (author_id, episode_id)
);

CREATE TABLE "radio_state"
(
    "author_id"         TEXT REFERENCES "author" ("id") ON DELETE CASCADE,
    "episode_id"        TEXT REFERENCES "radio_episode" ("episode_id") ON DELETE CASCADE,
    "started_at"        TIMESTAMP NULL,
    "current_timestamp" BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY ("author_id", "episode_id", "started_at")
);

CREATE INDEX latest_state ON radio_state ("author_id", "started_at");