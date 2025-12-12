CREATE TABLE flag_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id UUID NOT NULL REFERENCES flags(id) ON DELETE CASCADE,
    env_id UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    value TEXT,
    enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(flag_id, env_id)
);

CREATE INDEX idx_flag_values_flag_id ON flag_values(flag_id);
CREATE INDEX idx_flag_values_env_id ON flag_values(env_id);
CREATE INDEX idx_flag_values_enabled ON flag_values(enabled);
