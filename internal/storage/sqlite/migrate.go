package storage

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/migrations
var fs embed.FS

func migrateDB(dbURL string) {
	d, err := iofs.New(fs, "sql/migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	_, dirty, err := m.Version()
	if dirty || err != nil {
		panic("Wrong version")
	}

	err, err2 := m.Close()
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err2)
	}
}
