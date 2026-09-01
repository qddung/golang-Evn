package bookmark_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBookmarksByUserId(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		limit       int
		offset      int
		sort        string
		wantTotal   int64
		wantLength  int
		expectError func(t *testing.T, err error)
	}{
		{
			name:       "success returns bookmarks for user",
			userID:     "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			wantTotal:  1,
			wantLength: 1,
			expectError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:       "success with pagination and sort",
			userID:     "4e90220a-51f6-49e4-bc0e-44e2f321475a",
			limit:      10,
			offset:     0,
			sort:       "created_at desc",
			wantTotal:  1,
			wantLength: 1,
			expectError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:       "empty result for unknown user",
			userID:     "00000000-0000-0000-0000-000000000000",
			wantTotal:  0,
			wantLength: 0,
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

			bookmarks, total, err := repo.GetBookmarksByUserId(ctx, tc.userID, tc.limit, tc.offset, tc.sort)
			tc.expectError(t, err)
			assert.EqualValues(t, tc.wantTotal, total)
			assert.Len(t, bookmarks, tc.wantLength)
			for _, bookmark := range bookmarks {
				assert.Equal(t, tc.userID, bookmark.UserId)
				assert.NotNil(t, bookmark)
			}
		})
	}
}

