-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
   id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
   login VARCHAR(50) NOT NULL UNIQUE,
   password TEXT NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
