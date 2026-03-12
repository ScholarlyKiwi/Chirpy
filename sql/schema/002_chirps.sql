-- +goose Up
CREATE TABLE chirp (
    id          uuid,
    created_at  TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    body        TEXT NOT NULL,
    user_id     uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE chirp;