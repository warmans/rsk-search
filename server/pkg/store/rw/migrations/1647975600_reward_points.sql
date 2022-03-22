ALTER TABLE author_reward
    ADD COLUMN points_spent DECIMAL NOT NULL DEFAULT 0;

ALTER TABLE author_reward
    DROP COLUMN threshold;
