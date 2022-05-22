-- +goose Up
-- +goose StatementBegin
CREATE TABLE secrets (
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    secret_data BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
-- +goose StatementEnd
