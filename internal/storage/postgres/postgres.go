package postgres

import (
	"database/sql"
	"fmt"
	"golang-websocket-chat/internal/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Storage struct {
	db                 *bun.DB
	roomsRepository    *roomsRepository
	usersRepository    *usersRepository
	messagesRepository *messagesRepository
}

func New(conf config.Postgres) (*Storage, error) {
	const op = "storage.postgres.New"
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(conf.Address+":"+conf.Port),
		pgdriver.WithPassword(conf.Password),
		pgdriver.WithUser(conf.User),
		pgdriver.WithDatabase(conf.Db),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	db.RegisterModel((*UserToRoom)(nil))
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Rooms() *roomsRepository {
	if s.roomsRepository != nil {
		return s.roomsRepository
	}
	s.roomsRepository = &roomsRepository{
		db: s.db,
	}
	return s.roomsRepository
}

func (s *Storage) Users() *usersRepository {
	if s.usersRepository != nil {
		return s.usersRepository
	}
	s.usersRepository = &usersRepository{
		db: s.db,
	}
	return s.usersRepository
}

func (s *Storage) Messages() *messagesRepository {
	if s.messagesRepository != nil {
		return s.messagesRepository
	}
	s.messagesRepository = &messagesRepository{
		db: s.db,
	}
	return s.messagesRepository
}
