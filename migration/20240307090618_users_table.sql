-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  name text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
