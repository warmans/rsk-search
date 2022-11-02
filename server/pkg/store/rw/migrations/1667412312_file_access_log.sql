CREATE TABLE "media_access_log"
(
    epid               TEXT,
    media_type         TEXT,
    time_bucket        TIMESTAMP,
    num_times_accessed INTEGER,
    total_bytes        INTEGER,
    PRIMARY KEY (epid, media_type, time_bucket)
);
