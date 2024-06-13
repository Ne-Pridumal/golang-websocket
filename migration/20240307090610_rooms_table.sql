-- +goose Up
-- +goose StatementBegin
CREATE TABLE rooms (
  id SERIAL PRIMARY KEY,
  name text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
