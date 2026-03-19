INSERT INTO oauth_tokens (id, access_token, refresh_token, token_type, expiry_unix, updated_at_unix)
VALUES (1, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    access_token = excluded.access_token,
    refresh_token = excluded.refresh_token,
    token_type = excluded.token_type,
    expiry_unix = excluded.expiry_unix,
    updated_at_unix = excluded.updated_at_unix;
