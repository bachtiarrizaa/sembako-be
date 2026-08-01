CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(150) NOT NULL,
    email         VARCHAR(150) NOT NULL,
    username      VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    role_id       UUID NOT NULL REFERENCES roles(id),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_users_username ON users (username) WHERE deleted_at IS NULL AND username IS NOT NULL;
CREATE INDEX idx_users_deleted_at ON users (deleted_at);
CREATE INDEX idx_users_role_id ON users (role_id);