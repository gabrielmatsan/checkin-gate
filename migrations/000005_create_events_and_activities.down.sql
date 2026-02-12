DROP TRIGGER IF EXISTS update_activities_updated_at ON activities;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;

-- REMOVER INDEX
DROP INDEX IF EXISTS idx_activities_event_id;

DROP TABLE IF EXISTS activities;
DROP TABLE IF EXISTS events;