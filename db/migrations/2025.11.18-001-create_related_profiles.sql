CREATE TABLE IF NOT EXISTS related_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    real_profile_id UUID NOT NULL REFERENCES user_profiles(id),
    anon_profile_id UUID UNIQUE NOT NULL REFERENCES user_profiles(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
