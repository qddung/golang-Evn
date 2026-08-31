package bookmark_repository

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/test/data/fixture"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type inputCreateBookmark struct {
	userId      string
	url         string
	description string
	code        string
}

func (input *inputCreateBookmark) AssertCompareBookmark(t *testing.T, bookmark *entity.Bookmark, expected bool) {
	result := input.userId == bookmark.UserId &&
		input.url == bookmark.Url &&
		input.description == bookmark.Description &&
		input.code == bookmark.Code
	assert.Equal(t, expected, result)
}

func SetUpDB(t *testing.T) *gorm.DB {
	return fixture.NewFixture(t, fixture.NewBookmarkTestCase(t))
}

func TestCreateBookmark(t *testing.T) {
	testCase := []struct {
		name        string
		input       *inputCreateBookmark
		expectError func(t *testing.T, err error)
	}{
		{
			name: "success",
			input: &inputCreateBookmark{
				userId:      "test",
				url:         "test",
				description: "test",
				code:        "test",
			},
			expectError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			db := SetUpDB(t)
			repo := NewBookmarkRepository(db)
			bookmark, err := repo.CreateBookmark(ctx, tc.input.userId, tc.input.url, tc.input.description, tc.input.code)
			tc.expectError(t, err)
			if bookmark != nil {
				tc.input.AssertCompareBookmark(t, bookmark, true)
			}
		})
	}
}
