package postgres_test

import (
	"context"
	"golang-websocket-chat/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoomsRepository_Create(t *testing.T) {
	pg, d := testPostgres(t)
	ctx := context.Background()
	defer d("rooms")
	room := &models.Room{
		ID:   1,
		Name: "test",
	}
	err := pg.Rooms().Create(ctx, room)
	assert.NoError(t, err)
}

func TestRoomsRepository_Delete(t *testing.T) {
	id := 23423
	pg, d := testPostgres(t)
	ctx := context.Background()
	defer d("rooms")

	room := &models.Room{
		ID:   id,
		Name: "string",
	}

	pg.Rooms().Create(ctx, room)
	err := pg.Rooms().Delete(ctx, id)

	assert.NoError(t, err)
}

func TestRoomsRepository_GetById(t *testing.T) {
	id := 3423
	pg, d := testPostgres(t)
	ctx := context.Background()
	defer d("rooms", "messages")

	message := &models.Message{
		ID:      2342,
		RoomId:  id,
		Date:    time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		Content: "some content",
	}

	room := &models.Room{
		ID:       id,
		Name:     "string",
		Messages: []models.Message{*message},
	}

	err := pg.Rooms().Create(ctx, room)
	assert.NoError(t, err)

	err = pg.Messages().Create(ctx, message)
	assert.NoError(t, err)

	r, err := pg.Rooms().GetById(ctx, id)
	assert.NoError(t, err)
	assert.Exactly(t, room, r)
}

func TestRoomsRepository_AddUser(t *testing.T) {
	roomId := 32423
	userId := 12341234

	user := &models.User{
		ID:   userId,
		Name: "user",
	}

	room := &models.Room{
		ID:   roomId,
		Name: "test",
	}

	var usersSlice = []models.User{
		{
			ID:   userId,
			Name: "user",
		},
	}

	idlRoom := &models.Room{
		ID:    roomId,
		Name:  "test",
		Users: usersSlice,
	}

	pg, d := testPostgres(t)
	defer d("users", "rooms", "rooms_users")

	ctx := context.Background()

	pg.Rooms().Create(ctx, room)
	pg.Users().Create(ctx, user)

	err := pg.Rooms().AddUser(ctx, roomId, userId)
	assert.NoError(t, err)

	nRoom, err := pg.Rooms().GetById(ctx, roomId)
	assert.NoError(t, err)
	assert.Exactly(t, idlRoom, nRoom)
}
