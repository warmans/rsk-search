ALTER TABLE tscript_chunk
    ALTER COLUMN start_second TYPE bigint,
    ALTER COLUMN end_second TYPE bigint;


UPDATE tscript_chunk SET start_second=start_second*1000000000, end_second=greatest(end_second*1000000000, -1);