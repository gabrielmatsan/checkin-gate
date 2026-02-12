CREATE TABLE IF NOT EXISTS events (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  allowed_domains JSONB,
  description VARCHAR(500),
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS activities (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  event_id VARCHAR(36) NOT NULL REFERENCES events(id),
  description VARCHAR(500),
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

-- index para a coluna event_id
CREATE INDEX IF NOT EXISTS idx_activities_event_id ON activities (event_id);

-- trigger para atualizar a coluna updated_at
CREATE TRIGGER update_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- trigger para atualizar a coluna updated_at
CREATE TRIGGER update_activities_updated_at
BEFORE UPDATE ON activities
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


