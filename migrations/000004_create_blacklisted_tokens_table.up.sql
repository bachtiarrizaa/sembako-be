CREATE TABLE blacklisted_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash  VARCHAR(64) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_blacklisted_tokens_token_hash ON blacklisted_tokens (token_hash);
CREATE INDEX idx_blacklisted_tokens_expires_at ON blacklisted_tokens (expires_at);
