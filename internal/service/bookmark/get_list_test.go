package bookmark_service

import (
	"context"
	"testing"
	"time"

	"github.com/homework/lab/internal/models/base"
	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/models/entity"
	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	helpers_mocks "github.com/homework/lab/pkg/helpers/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_GetBookmarks(t *testing.T) {
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator)
		query        *bookmark_model.GetBookmarksQuery
		expectedFunc func(t *testing.T, res *api.PaginatedResponse[bookmark_model.BookmarkInfo], err error)
	}{
		{
			name: "success",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				repo.On("GetBookmarksByUserId", ctx, "user-1", 10, 0, "created_at desc").Return([]*entity.Bookmark{
					{
						Base:        base.Base{Id: "b-1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Url:         "https://example.com",
						Description: "demo bookmark",
						Code:        "ABC123",
						UserId:      "user-1",
					},
				}, int64(1), nil)
				return repo, keyGen
			},
			query: &bookmark_model.GetBookmarksQuery{Page: 1, Limit: 10, Sort: "created_at desc"},
			expectedFunc: func(t *testing.T, res *api.PaginatedResponse[bookmark_model.BookmarkInfo], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res.Data, 1)
				assert.EqualValues(t, 1, res.Pagination.Total)
				assert.EqualValues(t, 1, res.Pagination.Page)
				assert.EqualValues(t, 10, res.Pagination.Limit)
			},
		},
		{
			name: "repository error",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				repo.On("GetBookmarksByUserId", ctx, "user-1", 10, 0, "created_at desc").Return([]*entity.Bookmark(nil), int64(0), assert.AnError)
				return repo, keyGen
			},
			query: &bookmark_model.GetBookmarksQuery{Page: 1, Limit: 10, Sort: "created_at desc"},
			expectedFunc: func(t *testing.T, res *api.PaginatedResponse[bookmark_model.BookmarkInfo], err error) {
				assert.Equal(t, assert.AnError, err)
				assert.Nil(t, res)
			},
		},
		{
			name: "default page and limit",
			setupRepo: func(ctx context.Context) (*bookmark_mocks.BookmarkRepository, *helpers_mocks.KeyGenerator) {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				keyGen := helpers_mocks.NewKeyGenerator(t)
				repo.On("GetBookmarksByUserId", ctx, "user-1", 10, 0, "").Return([]*entity.Bookmark{}, int64(0), nil)
				return repo, keyGen
			},
			query: &bookmark_model.GetBookmarksQuery{},
			expectedFunc: func(t *testing.T, res *api.PaginatedResponse[bookmark_model.BookmarkInfo], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.EqualValues(t, 1, res.Pagination.Page)
				assert.EqualValues(t, 10, res.Pagination.Limit)
				assert.Empty(t, res.Data)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			repo, keyGen := tc.setupRepo(ctx)
			service := NewBookmarkService(repo, keyGen)
			res, err := service.GetBookmarks(ctx, "user-1", tc.query)
			tc.expectedFunc(t, res, err)
		})
	}
}
