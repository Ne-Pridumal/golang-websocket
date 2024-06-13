package postgres_test

import (
	"context"
	"golang-websocket-chat/internal/storage/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessagesRepository_Create(t *testing.T) {
	pg, d := testPostgres(t)
	ctx := context.Background()
	defer d("rooms", "messages")
	room := &postgres.Room{
		ID:   1,
		Name: "test",
	}

	err := pg.Rooms().Create(ctx, room)
	assert.NoError(t, err)

	message := &postgres.Message{
		ID:      2342,
		RoomId:  1,
		Date:    time.Now(),
		Content: "some content",
	}
	err = pg.Messages().Create(ctx, message)

	assert.NoError(t, err)
}

func TestMessageRepository_Delete(t *testing.T) {
	pg, d := testPostgres(t)
	ctx := context.Background()
	defer d("rooms", "messages")

	room := &postgres.Room{
		ID:   1,
		Name: "test",
	}

	err := pg.Rooms().Create(ctx, room)
	assert.NoError(t, err)

	message := &postgres.Message{
		ID:      2342,
		RoomId:  1,
		Date:    time.Now(),
		Content: "some content",
	}

	err = pg.Messages().Create(ctx, message)
	assert.NoError(t, err)

	err = pg.Messages().Delete(ctx, 2342)
	assert.NoError(t, err)
}
