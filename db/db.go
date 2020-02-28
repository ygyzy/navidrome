package db

import (
	"database/sql"
	"os"
	"sync"

	"github.com/deluan/navidrome/conf"
	_ "github.com/deluan/navidrome/db/migration"
	"github.com/deluan/navidrome/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var (
	Driver = "sqlite3"
	Path   string
)

var (
	once sync.Once
	db   *sql.DB
)

func Db() *sql.DB {
	once.Do(func() {
		var err error
		Path = conf.Server.DbPath
		if Path == ":memory:" {
			Path = "file::memory:?cache=shared"
			conf.Server.DbPath = Path
		}
		log.Debug("Opening DataBase", "dbPath", Path, "driver", Driver)
		db, err = sql.Open(Driver, Path)
		if err != nil {
			panic(err)
		}
	})
	return db
}

func EnsureLatestVersion() {
	db := Db()

	err := goose.SetDialect(Driver)
	if err != nil {
		log.Error("Invalid DB driver", "driver", Driver, err)
		os.Exit(1)
	}
	err = goose.Run("up", db, "./")
	if err != nil {
		log.Error("Failed to apply new migrations", err)
		os.Exit(1)
	}
}
