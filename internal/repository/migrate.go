package repository

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"

	"firstmod/migrations"
)

func (db *DB) Migrate() error {
	db.log.Debug("running migration")
	files, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		db.log.Error("failed to load migration files", "error", err)
		return err
	}
	db.log.Debug("migration files loaded successfully")

	sqlDB := stdlib.OpenDBFromPool(db.conn)
	defer sqlDB.Close()

	driver, err := pgx.WithInstance(sqlDB, &pgx.Config{})
	if err != nil {
		db.log.Error("failed to create pgx driver for migrations", "error", err)
		return err
	}
	m, err := migrate.NewWithInstance("iofs", files, "pgx", driver)
	if err != nil {
		db.log.Error("failed to initialize migrations", "error", err)
		return err
	}

	err = m.Up()

	if err != nil {
		if err != migrate.ErrNoChange {
			db.log.Error("migration failed", "error", err)
			return err
		}
		db.log.Debug("migration did not change anything")
	}

	db.log.Debug("migration finished")
	return nil
}
