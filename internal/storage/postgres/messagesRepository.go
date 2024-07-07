package postgres

import (
	"context"
	"golang-websocket-chat/internal/lib"
	"golang-websocket-chat/internal/models"

	"github.com/uptrace/bun"
)

type messagesRepository struct {
	db *bun.DB
}

func (r *messagesRepository) Create(ctx context.Context, message *models.Message) error {
	const op = "storage.postgres.messagesRep.Create"

	_, err := r.db.NewInsert().Model(message).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *messagesRepository) Delete(ctx context.Context, id int) error {
	const op = "storage.postgres.messagesRepo.Delete"

	_, err := r.db.NewDelete().Model(&models.Message{}).Where("id = ?", id).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}
