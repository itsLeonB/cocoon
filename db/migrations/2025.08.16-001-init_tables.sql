CREATE TYPE friendship_type AS ENUM (
    'REAL',
    'ANON'
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
COMMENT ON COLUMN user_profiles.user_id IS 'Nullable. Can be NULL for peers who do not have an account in the app';

CREATE TABLE IF NOT EXISTS friendships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id1 UUID NOT NULL REFERENCES user_profiles(id),
    profile_id2 UUID NOT NULL REFERENCES user_profiles(id),
    type friendship_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_friendship UNIQUE (profile_id1, profile_id2),
    CONSTRAINT profile_order CHECK (profile_id1 < profile_id2)
);

CREATE INDEX IF NOT EXISTS user_profiles_user_id_idx ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS user_profiles_name_idx ON user_profiles(name);
CREATE INDEX IF NOT EXISTS friendships_profile_id1_idx ON friendships(profile_id1);
CREATE INDEX IF NOT EXISTS friendships_profile_id2_idx ON friendships(profile_id2);
CREATE INDEX IF NOT EXISTS friendships_type_idx ON friendships(type);
