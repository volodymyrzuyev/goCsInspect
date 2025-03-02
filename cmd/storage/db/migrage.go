package db

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs" // Ensure import
	"github.com/volodymyrzuyev/goCsInspect/cmd/logger"
)

//go:embed sqlc/sql/*.sql
var migrations embed.FS

func migrateUp(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		logger.ERROR.Printf("Unable to create a DB driver. Err: %v", err)
		return err
	}

	sourceInstance, err := iofs.New(migrations, "sqlc/sql")
	if err != nil {
		logger.ERROR.Printf("Can not make source index for the migration folder %v", err)
		return err
	}

	m, err := migrate.NewWithInstance(
		"sql",
		sourceInstance,
		"sqlite3",
		driver,
	)
	if err != nil {
		logger.ERROR.Printf("Can't create new migration instance. Err: %v", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.ERROR.Printf("Unable to migrate up. Err: %v", err)
		return err
	}

	return nil
}

