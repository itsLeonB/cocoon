-- Enable trigram extension once
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create search table if it doesn't exist
CREATE TABLE IF NOT EXISTS profile_names (
    id UUID PRIMARY KEY,
    name TEXT
);

-- Seed data for existing valid profiles
INSERT INTO profile_names (id, name)
SELECT id, name
FROM user_profiles
WHERE user_id IS NOT NULL
  AND user_id <> '00000000-0000-0000-0000-000000000000'
ON CONFLICT (id) DO NOTHING;

-- Create trigram index for fuzzy search
CREATE INDEX IF NOT EXISTS idx_profile_names_name_trgm
ON profile_names
USING GIN (name gin_trgm_ops);

-- Define trigger sync function
CREATE OR REPLACE FUNCTION sync_profile_names()
RETURNS TRIGGER AS $$
BEGIN
  IF (TG_OP = 'INSERT') THEN
    IF NEW.user_id IS NOT NULL AND NEW.user_id <> '00000000-0000-0000-0000-000000000000' THEN
      INSERT INTO profile_names (id, name)
      VALUES (NEW.id, NEW.name);
    END IF;

  ELSIF (TG_OP = 'UPDATE') THEN
    IF NEW.user_id IS NOT NULL AND NEW.user_id <> '00000000-0000-0000-0000-000000000000' THEN
      INSERT INTO profile_names (id, name)
      VALUES (NEW.id, NEW.name)
      ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;
    ELSE
      DELETE FROM profile_names WHERE id = OLD.id;
    END IF;

  ELSIF (TG_OP = 'DELETE') THEN
    DELETE FROM profile_names WHERE id = OLD.id;
  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger only if missing
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'trg_sync_profile_names'
  ) THEN
    CREATE TRIGGER trg_sync_profile_names
    AFTER INSERT OR UPDATE OR DELETE ON user_profiles
    FOR EACH ROW
    EXECUTE FUNCTION sync_profile_names();
  END IF;
END$$;
