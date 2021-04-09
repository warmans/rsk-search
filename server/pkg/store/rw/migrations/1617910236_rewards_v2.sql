ALTER TABLE "author_reward"
    RENAME COLUMN confirmed to claimed;

ALTER TABLE "author_reward"
    ADD COLUMN claim_kind              TEXT,
    ADD COLUMN claim_value             NUMERIC   NULL,
    ADD COLUMN claim_value_currency    TEXT      NULL,
    ADD COLUMN claim_description       TEXT      NULL,
    ADD COLUMN claim_at                TIMESTAMP NULL,
    ADD COLUMN claim_confirmation_code TEXT      NULL;
