ALTER TABLE "media_access_log" ADD COLUMN total_mib BIGINT;

UPDATE media_access_log SET total_mib = total_bytes / 1048576;

