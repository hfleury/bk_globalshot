DROP INDEX IF EXISTS idx_auth_tokens_user_id;
DROP INDEX IF EXISTS idx_auth_tokens_expires;
DROP INDEX IF EXISTS idx_auth_tokens_revoked;

DROP TABLE IF EXISTS auth_tokens;