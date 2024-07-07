package postgres

import (
	"context"
	"fmt"
	"golang-websocket-chat/internal/lib"
	"golang-websocket-chat/internal/models"

	"github.com/uptrace/bun"
)

type usersRepository struct {
	db *bun.DB
}

func (r *usersRepository) Create(ctx context.Context, user *models.User) error {
	const op = "storage.postgres.userRep.Create"

	_, err := r.db.NewInsert().Model(user).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *usersRepository) Delete(ctx context.Context, id int) error {
	const op = "storage.postgres.userRep.Delete"

	_, err := r.db.NewDelete().Model(&models.User{}).Where("id = ?", id).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *usersRepository) GetById(ctx context.Context, id int) (*models.User, error) {
	const op = "storage.postgres.userRep.GetById"

	usr := new(models.User)
	err := r.db.NewSelect().Model(usr).Relation("Rooms").Where("id = ?", id).Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usr, nil
}
