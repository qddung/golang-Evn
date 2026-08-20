package sqldb

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func CatchDBError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// BuildMigrate return migrate.Migrate
func BuildMigrate(db *gorm.DB, migrationPath string) *migrate.Migrate {
	sqlDb, err := db.DB()
	CatchDBError(err)
	pgDriver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	CatchDBError(err)
	migrator, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", migrationPath),
		db.Name(), pgDriver)
	CatchDBError(err)
	return migrator
}

// MigrateUp and MigrateDown
func MigrateUp(migrator *migrate.Migrate, step ...int) error {
	defaultStep := 1
	if len(step) > 0 {
		defaultStep = step[0]
	}
	return migrator.Steps(defaultStep)
}

// MigrateUp and MigrateDown
func MigrateDown(migrator *migrate.Migrate, step ...int) error {
	defaultStep := -1
	if len(step) > 0 && step[0] < -1 {
		defaultStep = step[0]
	}
	return migrator.Steps(defaultStep)
}
