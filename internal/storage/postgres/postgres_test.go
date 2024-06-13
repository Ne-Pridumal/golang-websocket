package postgres_test

import (
	"database/sql"
	"fmt"
	"golang-websocket-chat/internal/config"
	"golang-websocket-chat/internal/storage/postgres"
	"strings"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func testPostgres(t *testing.T) (*postgres.Storage, func(...string)) {
	t.Helper()
	conf := config.Postgres{
		User:     "root",
		Password: "root",
		Address:  "localhost",
		Port:     "5432",
		Ssl:      false,
		Db:       "test",
	}

	s, err := postgres.New(conf)

	if err != nil {
		panic(err)
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(conf.Address+":"+conf.Port),
		pgdriver.WithPassword(conf.Password),
		pgdriver.WithUser(conf.User),
		pgdriver.WithDatabase(conf.Db),
		pgdriver.WithInsecure(true),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	return s, func(tables ...string) {
		if len(tables) > 0 {
			str := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ","))
			db.Exec(str)
		}
		db.Close()
	}
}
