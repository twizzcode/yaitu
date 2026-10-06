CREATE TABLE auth_codes (
    code_hash BYTEA PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_host TEXT NOT NULL,
    return_path TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (target_host = LOWER(target_host)),
    CHECK (target_host NOT LIKE '%/%'),
    CHECK (target_host NOT LIKE '%:%'),
    CHECK (target_host NOT LIKE '.%'),
    CHECK (target_host NOT LIKE '%.'),
    CHECK (target_host <> ''),
    CHECK (return_path LIKE '/%'),
    CHECK (return_path NOT LIKE '//%')
);

CREATE INDEX auth_codes_user_id_idx
ON auth_codes(user_id);

CREATE INDEX auth_codes_expires_at_idx
ON auth_codes(expires_at);