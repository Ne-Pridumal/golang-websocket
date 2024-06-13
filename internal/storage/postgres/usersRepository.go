package postgres

import (
	"context"
	"fmt"
	"golang-websocket-chat/internal/lib"

	"github.com/uptrace/bun"
)

type usersRepository struct {
	db *bun.DB
}

func (r *usersRepository) Create(ctx context.Context, user *User) error {
	const op = "storage.postgres.userRep.Create"

	_, err := r.db.NewInsert().Model(user).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *usersRepository) Delete(ctx context.Context, id int) error {
	const op = "storage.postgres.userRep.Delete"

	_, err := r.db.NewDelete().Model(&User{}).Where("id = ?", id).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *usersRepository) GetById(ctx context.Context, id int) (*User, error) {
	const op = "storage.postgres.userRep.GetById"

	usr := new(User)
	err := r.db.NewSelect().Model(usr).Relation("Rooms").Where("id = ?", id).Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usr, nil
}
