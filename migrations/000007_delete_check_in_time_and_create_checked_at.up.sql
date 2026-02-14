ALTER TABLE check_ins DROP COLUMN check_in_time;
ALTER TABLE check_ins ADD COLUMN checked_at TIMESTAMPTZ;