-- +goose up
CREATE TABLE chirps (
    id UUID Primary Key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body string NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE chirps;