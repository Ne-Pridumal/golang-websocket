-- +goose Up
-- +goose StatementBegin
CREATE TABLE rooms_users (
  room_id int REFERENCES rooms(id) ON UPDATE CASCADE ON DELETE CASCADE,
  user_id int REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
