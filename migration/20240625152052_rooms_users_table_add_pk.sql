-- +goose Up
-- +goose StatementBegin
ALTER TABLE rooms_users ADD PRIMARY KEY(room_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rooms_users;
-- +goose StatementEnd
