package storage

import (
	"embed"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/migrations
var fs embed.FS

func migrateDB(dbURL string, l *slog.Logger) error {
	d, err := iofs.New(fs, "sql/migrations")
	if err != nil {
		l.Error("error getting migrations", "error", err)
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		l.Error("error getting migrations", "error", err)
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		l.Error("error applying migrations", "error", err)
		return err
	}

	_, dirty, err := m.Version()
	if dirty || err != nil {
		l.Error("dirty storage version", "error", err)
		return fmt.Errorf("dirty storage version")
	}

	err, err2 := m.Close()
	if err != nil {
		l.Error("error closing", "error", err)
		return err
	}
	if err2 != nil {
		l.Error("error closing", "error", err)
		return err
	}

	return nil
}
