-- These tables allow changes/chunk contributions to be linked back to rewards so users can see how many points
-- were awarded for each change.
CREATE TABLE "author_contribution_transcript_change"
(
    author_contribution_id TEXT REFERENCES author_contribution (id),
    transcript_change_id   TEXT REFERENCES transcript_change (id)
);

CREATE TABLE "author_contribution_tscript_contribution"
(
    author_contribution_id  TEXT REFERENCES author_contribution (id),
    tscript_contribution_id TEXT REFERENCES tscript_contribution (id)
);
