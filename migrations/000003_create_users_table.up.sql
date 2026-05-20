CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL    PRIMARY KEY,
    username   VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    name       VARCHAR(255),
    role       VARCHAR(50)  NOT NULL DEFAULT 'user',
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
 
CREATE INDEX IF NOT EXISTS idx_users_username   ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
 
INSERT INTO users (username, password, name, role, created_by, updated_by) VALUES
    ('admin', 'admin123', 'Administrator', 'admin', 'system', 'system'),
    ('user',  'user123',  'User',          'user',  'system', 'system')
ON CONFLICT (username) DO NOTHING;