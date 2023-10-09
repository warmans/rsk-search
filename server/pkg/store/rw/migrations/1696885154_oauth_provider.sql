ALTER TABLE author
    ADD COLUMN oauth_provider TEXT DEFAULT 'reddit';

ALTER TABLE author
    DROP CONSTRAINT "author_name_key";

CREATE UNIQUE INDEX author_name_oauth_provider ON author (name, oauth_provider);
