ALTER TABLE check_ins ADD COLUMN check_in_time TIMESTAMPTZ;
ALTER TABLE check_ins DROP COLUMN checked_at;