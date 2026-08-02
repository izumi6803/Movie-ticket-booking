-- Allow an email to be reused after its previous account is soft-deleted.
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
DROP INDEX IF EXISTS idx_users_email;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
    ON users (LOWER(email))
    WHERE deleted_at IS NULL;
