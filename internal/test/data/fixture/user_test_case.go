package fixture

import (
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"gorm.io/gorm"
)

type user_test_case struct {
	base
}

func (u *user_test_case) Migrate() error {

	return u.db.AutoMigrate(&entity.User{})
}

func (u *user_test_case) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []*entity.User{
		{DisplayName: "brad", Email: "abc@gmail.com", Password: "123456", UserName: "acd"},
		{DisplayName: "brad2", Email: "abc2@gmail.com", Password: "123456", UserName: "acd2"},
		{DisplayName: "brad3", Email: "abc3@gmail.com", Password: "123456", UserName: "acd3"},
	}
	return db.CreateInBatches(users, 10).Error
}

func NewUserTestCase(t *testing.T) Fixture {
	return &user_test_case{}
}
