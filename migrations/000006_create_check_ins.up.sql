CREATE TABLE IF NOT EXISTS check_ins (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL REFERENCES users(id),
  activity_id VARCHAR(36) NOT NULL REFERENCES activities(id),
  check_in_time TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE INDEX IF NOT EXISTS idx_check_ins_user_id ON check_ins (user_id);
CREATE INDEX IF NOT EXISTS idx_check_ins_activity_id ON check_ins (activity_id);




