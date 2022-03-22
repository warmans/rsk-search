ALTER TABLE author_contribution
    ADD COLUMN points_spent_v2 DECIMAL NOT NULL DEFAULT 0;

UPDATE author_contribution
SET points_spent_v2 = points
WHERE points_spent = true;

ALTER TABLE author_contribution
    DROP COLUMN points_spent;

ALTER TABLE author_contribution
    RENAME COLUMN points_spent_v2 TO points_spent;

