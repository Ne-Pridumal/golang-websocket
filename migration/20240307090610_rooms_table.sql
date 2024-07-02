-- +goose Up
-- +goose StatementBegin
CREATE TABLE rooms (
  id SERIAL PRIMARY KEY,
  name text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rooms;
-- +goose StatementEnd
