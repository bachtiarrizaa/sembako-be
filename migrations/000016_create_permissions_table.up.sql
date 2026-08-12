CREATE TABLE permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id   UUID REFERENCES permissions(id) ON DELETE CASCADE,
    type        VARCHAR(20) NOT NULL DEFAULT 'action' CHECK (type IN ('menu', 'action')),
    path        VARCHAR(255),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_permissions_parent ON permissions(parent_id);
