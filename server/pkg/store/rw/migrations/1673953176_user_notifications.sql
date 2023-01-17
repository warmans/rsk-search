CREATE TABLE "author_notification"
(
    id                TEXT      NOT NULL,
    author_id         TEXT REFERENCES author(id) ON DELETE CASCADE,
    kind              TEXT      NOT NULL,
    message           TEXT      NOT NULL,
    click_through_url TEXT NULL,
    read_at           TIMESTAMP NULL,
    created_at        TIMESTAMP NOT NULL,
    PRIMARY KEY (id)
);
