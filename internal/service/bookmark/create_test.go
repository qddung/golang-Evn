package bookmark_service

import (
	"context"
	"testing"
	"time"

	"github.com/homework/lab/internal/models/base"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/models/entity"
	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	helpers_mocks "github.com/homework/lab/pkg/helpers/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestService_NewBookmark(t *testing.T) {
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator)
		input        *bookmark_model.NewBookmarkRequest
		expectedFunc func(t *testing.T, info *bookmark_model.BookmarkInfo, err error)
	}{
		{
			name: "success",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				keyGen.On("GenerateRandomCode", 10).Return("ABC123XYZ9")
				repo.On("CreateBookmark", ctx, "user-1", "https://example.com", "demo bookmark", "ABC123XYZ9").Return(&entity.Bookmark{
					Base:        base.Base{Id: "b-1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Url:         "https://example.com",
					Description: "demo bookmark",
					Code:        "ABC123XYZ9",
					UserId:      "user-1",
				}, nil)
				return repo, keyGen
			},
			input: &bookmark_model.NewBookmarkRequest{
				Url:         "https://example.com",
				Description: "demo bookmark",
			},
			expectedFunc: func(t *testing.T, info *bookmark_model.BookmarkInfo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, info)
				assert.Equal(t, "https://example.com", info.Url)
				assert.Equal(t, "demo bookmark", info.Description)
				assert.Equal(t, "ABC123XYZ9", info.Code)
			},
		},
		{
			name: "repository error",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				keyGen.On("GenerateRandomCode", 10).Return("ABC123XYZ9")
				repo.On("CreateBookmark", ctx, "user-1", "https://example.com", "demo bookmark", "ABC123XYZ9").Return(nil, assert.AnError)
				return repo, keyGen
			},
			input: &bookmark_model.NewBookmarkRequest{
				Url:         "https://example.com",
				Description: "demo bookmark",
			},
			expectedFunc: func(t *testing.T, info *bookmark_model.BookmarkInfo, err error) {
				assert.Equal(t, response.InternalError, err)
				assert.Nil(t, info)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, keyGen := tc.setupRepo(ctx)
			service := NewBookmarkService(repo, keyGen)
			info, err := service.NewBookmark(ctx, "user-1", tc.input)
			tc.expectedFunc(t, info, err)
		})
	}
}
