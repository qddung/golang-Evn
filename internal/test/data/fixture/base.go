package fixture

import (
	"testing"

	"github.com/homework/lab/pkg/sqldb"
	gorm "gorm.io/gorm"
)

type Fixture interface {
	SetDB(db *gorm.DB)
	DB() *gorm.DB
	// virtual method in base struct
	Migrate() error
	GenerateData() error
}

// abstract class - no instance create
type base struct {
	db *gorm.DB
}

func (b *base) DB() *gorm.DB {
	return b.db
}

func (b *base) SetDB(db *gorm.DB) {
	b.db = db
}

// setup db and data for test
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	db, err := sqldb.NewMiniPostgres(t)
	if err != nil {
		t.Fatalf("Failed to start DB: %v", err)
	}

	if fix == nil {
		return db
	}
	fix.SetDB(db)
	err = fix.Migrate()
	if err != nil {
		t.Fatalf("Failed to migrate data: %v", err)
	}

	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("Failed to generate data: %v", err)
	}
	return fix.DB()
}
