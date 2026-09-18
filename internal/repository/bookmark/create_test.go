package bookmark_repository

import (
	"context"
	"strings"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	"github.com/homework/lab/internal/test/data/fixture"
	base62_helper "github.com/homework/lab/pkg/helpers/base62"
	base62_lib "github.com/homework/lab/pkg/lib/base62"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type inputCreateBookmark struct {
	userId      string
	url         string
	description string
}

var prefix = "bookmark_url_"

func (input *inputCreateBookmark) AssertCompareBookmark(t *testing.T, bookmark *entity.Bookmark, expected bool) {
	result := input.userId == bookmark.UserId &&
		input.url == bookmark.Url &&
		input.description == bookmark.Description && strings.Contains(bookmark.Code, prefix)
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
			code_gen := base62_helper.New(base62_lib.NewStdEncoding())
			repo := NewBookmarkRepository(db, code_gen)
			bookmark, err := repo.CreateBookmark(ctx, tc.input.userId, tc.input.url, tc.input.description, prefix)
			tc.expectError(t, err)
			if bookmark != nil {
				tc.input.AssertCompareBookmark(t, bookmark, true)
			}
		})
	}
}
