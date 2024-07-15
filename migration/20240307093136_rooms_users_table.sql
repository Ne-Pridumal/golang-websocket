-- +goose Up
-- +goose StatementBegin
CREATE TABLE rooms_users (
  room_id int REFERENCES rooms(id) ON UPDATE CASCADE ON DELETE CASCADE,
  user_id int REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  PRIMARY KEY (room_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE room_users;
-- +goose StatementEnd
