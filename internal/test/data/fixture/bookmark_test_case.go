package fixture

import (
	"testing"

	"github.com/homework/lab/internal/models/entity"
	gorm "gorm.io/gorm"
)

type bookmark_test_case struct {
	user_test_case
}

func (u *bookmark_test_case) Migrate() error {
	u.user_test_case.Migrate()
	return u.db.AutoMigrate(&entity.Bookmark{})
}

func (u *bookmark_test_case) GenerateData() error {
	gen_user_error := u.user_test_case.GenerateData()
	if gen_user_error != nil {
		return gen_user_error
	}
	db := u.db.Session(&gorm.Session{SkipHooks: true})
	bookmarks := []*entity.Bookmark{
		{
			Base:        GetBaseEntity("28d51812-2c60-4896-853b-cfccd06d5243"),
			Url:         "https://github.com/bradfitz/gomemfs",
			UserId:      "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			Code:        "1234567890",
			Description: "Test 1",
		},
		{
			Base:        GetBaseEntity("28d51812-2c60-4896-853b-cfccd06d5244"),
			Url:         "https://github.com/bradfitz/gomemfs",
			UserId:      "4e90220a-51f6-49e4-bc0e-44e2f321476a",
			Code:        "1234567890",
			Description: "Test 1",
		},

		{
			Base:        GetBaseEntity("28d51812-2c60-4896-853b-cfccd06d5246"),
			Url:         "https://github.com/bradfitz/gomemfs",
			UserId:      "4e90220a-51f6-49e4-bc0e-44e2f321477a",
			Code:        "1234567890",
			Description: "Test 1",
		},
	}
	return db.CreateInBatches(bookmarks, 10).Error
}

func NewBookmarkTestCase(t *testing.T) Fixture {
	return &bookmark_test_case{}
}
