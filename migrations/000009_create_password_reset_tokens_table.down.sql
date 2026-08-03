DROP TRIGGER IF EXISTS trg_reset_tokens_updated_at ON password_reset_tokens;
DROP INDEX IF EXISTS idx_reset_tokens_hash;
DROP INDEX IF EXISTS idx_reset_tokens_user;
DROP TABLE IF EXISTS password_reset_tokens;
DROP FUNCTION IF EXISTS set_updated_at();
