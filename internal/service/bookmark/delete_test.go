package bookmark_service

import (
	"context"
	"testing"

	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestService_DeleteBookmark(t *testing.T) {
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context) *bookmark_mocks.BookmarkRepository
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupRepo: func(ctx context.Context) *bookmark_mocks.BookmarkRepository {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				repo.On("DeleteBookmark", ctx, "user-1", "bookmark-1").Return(nil)
				return repo
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "repository error",
			setupRepo: func(ctx context.Context) *bookmark_mocks.BookmarkRepository {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				repo.On("DeleteBookmark", ctx, "user-1", "bookmark-1").Return(gorm.ErrRecordNotFound)
				return repo
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Equal(t, response.NotFoundError, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo := tc.setupRepo(ctx)
			service := NewBookmarkService(repo)
			err := service.DeleteBookmark(ctx, "user-1", "bookmark-1")
			tc.expectedFunc(t, err)
		})
	}
}
