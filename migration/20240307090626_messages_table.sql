-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
  id SERIAL PRIMARY KEY,
  room_id int REFERENCES rooms(id) ON UPDATE CASCADE ON DELETE CASCADE, 
  content text not null,
  date timestamp not null,
  CONSTRAINT fk_room FOREIGN KEY(room_id) REFERENCES rooms(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE messages;
-- +goose StatementEnd
