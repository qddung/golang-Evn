package bookmark_service

import (
	"context"
	"testing"

	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	helpers_mocks "github.com/homework/lab/pkg/helpers/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestService_UpdateBookmark(t *testing.T) {
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator)
		input        *bookmark_model.UpdateBookmarkRequest
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name: "success",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				repo.On("UpdateBookmark", ctx, "user-1", "bookmark-1", "https://example.com/updated", "updated description").Return(nil)
				return repo, keyGen
			},
			input: &bookmark_model.UpdateBookmarkRequest{
				Url:         "https://example.com/updated",
				Description: "updated description",
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "repository error",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				repo.On("UpdateBookmark", ctx, "user-1", "bookmark-1", "https://example.com/updated", "updated description").Return(gorm.ErrRecordNotFound)
				return repo, keyGen
			},
			input: &bookmark_model.UpdateBookmarkRequest{
				Url:         "https://example.com/updated",
				Description: "updated description",
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
			repo, keyGen := tc.setupRepo(ctx)
			service := NewBookmarkService(repo, keyGen)
			err := service.UpdateBookmark(ctx, tc.input, "user-1", "bookmark-1")
			tc.expectedFunc(t, err)
		})
	}
}
