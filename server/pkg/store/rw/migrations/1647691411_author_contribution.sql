CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- In order to make it possible to remove completed transcript chunks from the DB,
-- the author contributions need to be tracked separately.
--

CREATE TABLE "author_contribution"
(
    id                TEXT PRIMARY KEY,
    author_id         TEXT REFERENCES author (id) ON DELETE CASCADE,
    epid              TEXT,
    contribution_type TEXT,
    points            DECIMAL,
    points_spent      BOOLEAN,
    created_at        TIMESTAMP
);

INSERT INTO author_contribution (id, author_id, epid, contribution_type, points, points_spent, created_at)
SELECT uuid_generate_v4(), co.author_id, replace(ch.tscript_id, 'ts-', 'ep-'), 'chunk', 1, true, co.created_at
FROM tscript_contribution co
LEFT JOIN tscript_chunk ch ON co.tscript_chunk_id = ch.id
WHERE co.state = 'approved';


INSERT INTO author_contribution (id, author_id, epid, contribution_type, points, points_spent, created_at)
SELECT uuid_generate_v4(), ch.author_id, ch.epid, 'change', 1, true, ch.created_at
FROM transcript_change ch
WHERE ch.state = 'approved';
