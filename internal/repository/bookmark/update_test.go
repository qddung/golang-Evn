package bookmark_repository

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateBookmark(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		bookmarkID  string
		url         string
		description string
		assertFunc  func(t *testing.T, db *gorm.DB)
		expectError func(t *testing.T, err error)
	}{
		{
			name:        "success updates bookmark",
			userID:      "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			bookmarkID:  "28d51812-2c60-4896-853b-cfccd06d5243",
			url:         "https://example.com/updated",
			description: "updated description",
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var bookmark entity.Bookmark
				err := db.Where("id = ?", "28d51812-2c60-4896-853b-cfccd06d5243").First(&bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, "https://example.com/updated", bookmark.Url)
				assert.Equal(t, "updated description", bookmark.Description)
			},
			expectError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:        "error when bookmark not found",
			userID:      "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			bookmarkID:  "00000000-0000-0000-0000-000000000000",
			url:         "https://example.com/invalid",
			description: "should not update",
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				err := db.Model(&entity.Bookmark{}).Count(&count).Error
				assert.NoError(t, err)
				assert.EqualValues(t, 3, count)
			},
			expectError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			db := SetUpDB(t)
			repo := NewBookmarkRepository(db)

			err := repo.UpdateBookmark(ctx, tc.userID, tc.bookmarkID, tc.url, tc.description)
			tc.expectError(t, err)
			if tc.assertFunc != nil {
				tc.assertFunc(t, db)
			}
		})
	}
}
