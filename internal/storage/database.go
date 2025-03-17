package storage

import (
	"errors"
	"github.com/xo/dburl"
	"log"
	"mailer/internal"

	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	migrateTargetImpl "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDatabase() *sql.DB {
	c, err := sql.Open("pgx", internal.Config.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return c
}

func parseDatabaseName() (string, error) {
	parsed, err := dburl.Parse(internal.Config.DatabaseURL)
	if err != nil {
		return "", err
	}
	return parsed.Path, nil
}

func (s *Storage) MigrateDatabase() {
	driver, err := migrateTargetImpl.WithInstance(s.Database, &migrateTargetImpl.Config{})
	if err != nil {
		log.Fatalf("Unable to migrate the database: %+v", err)
	}

	databaseName, err := parseDatabaseName()
	if err != nil {
		log.Fatalf("Unable to parse the database name: %+v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", databaseName, driver)
	if err != nil {
		log.Fatalf("Unable to migrate the database: %+v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Unable to migrate the database: %+v", err)
		return
	}
}
