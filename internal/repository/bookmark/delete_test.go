package bookmark_repository

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteBookmark(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		bookmarkID  string
		assertFunc  func(t *testing.T, db *gorm.DB)
		expectError func(t *testing.T, err error)
	}{
		{
			name:       "success deletes bookmark",
			userID:     "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			bookmarkID: "28d51812-2c60-4896-853b-cfccd06d5243",
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				err := db.Model(&entity.Bookmark{}).Where("id = ?", "28d51812-2c60-4896-853b-cfccd06d5243").Count(&count).Error
				assert.NoError(t, err)
				assert.EqualValues(t, 0, count)
			},
			expectError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:       "delete unknown bookmark succeeds without error",
			userID:     "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			bookmarkID: "00000000-0000-0000-0000-000000000000",
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				err := db.Model(&entity.Bookmark{}).Count(&count).Error
				assert.NoError(t, err)
				assert.EqualValues(t, 3, count)
			},
			expectError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			db := SetUpDB(t)
			repo := NewBookmarkRepository(db)

			err := repo.DeleteBookmark(ctx, tc.userID, tc.bookmarkID)
			tc.expectError(t, err)
			if tc.assertFunc != nil {
				tc.assertFunc(t, db)
			}
		})
	}
}
