-- +goose up
CREATE TABLE refresh_tokens (
    token TEXT Primary Key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP  

);

-- +goose down
DROP TABLE refresh_tokens;