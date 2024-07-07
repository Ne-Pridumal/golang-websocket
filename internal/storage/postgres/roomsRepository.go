package postgres

import (
	"context"
	"fmt"
	"golang-websocket-chat/internal/lib"
	"golang-websocket-chat/internal/models"

	"github.com/uptrace/bun"
)

type roomsRepository struct {
	db *bun.DB
}

func (r *roomsRepository) Create(ctx context.Context, room *models.Room) error {
	const op = "storage.postgres.roomsRep.Create"

	_, err := r.db.NewInsert().Model(room).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *roomsRepository) Delete(ctx context.Context, id int) error {
	const op = "storage.postgres.roomsRep.Delete"
	rm := new(models.Room)
	_, err := r.db.NewDelete().Model(rm).Where("id = ?", id).Exec(ctx)

	return lib.ErrWrapper(err, op)
}

func (r *roomsRepository) GetById(ctx context.Context, id int) (*models.Room, error) {
	const op = "storage.postgres.roomsRep.GetById"

	rm := new(models.Room)
	err := r.db.NewSelect().Model(rm).Relation("Users").Relation("Messages").Where("id = ?", id).Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rm, nil
}

func (r *roomsRepository) AddUser(ctx context.Context, roomId int, userId int) error {
	const op = "storage.postgres.roomsRep.AddUser"

	rm := &models.UserToRoom{
		RoomID: roomId,
		UserID: userId,
	}
	_, err := r.db.NewInsert().Model(rm).Returning("NULL").Exec(ctx)

	return lib.ErrWrapper(err, op)
}
