ALTER TABLE "author"
    ADD COLUMN "placeholder" BOOLEAN DEFAULT FALSE;

CREATE TABLE "episode_review"
(
    "author_id"  TEXT      NULL REFERENCES "author" ("id") ON DELETE CASCADE NOT NULL,
    "episode_id" TEXT      NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "rating"     DECIMAL   NOT NULL,
    "review"     TEXT      NULL,
    PRIMARY KEY ("author_id", "episode_id")
);
