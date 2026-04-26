ALTER TABLE outbox
DROP COLUMN attempts,
DROP COLUMN last_error,
DROP COLUMN next_retry_at;