package bookmark_service

import (
	"context"
	"testing"
	"time"

	"github.com/homework/lab/internal/models/base"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/models/entity"
	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestService_NewBookmark(t *testing.T) {
	testCases := []struct {
		name         string
		setupRepo    func(ctx context.Context) *bookmark_mocks.BookmarkRepository
		input        *bookmark_model.NewBookmarkRequest
		expectedFunc func(t *testing.T, info *bookmark_model.BookmarkInfo, err error)
	}{
		{
			name: "success",
			setupRepo: func(ctx context.Context) *bookmark_mocks.BookmarkRepository {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				repo.On("CreateBookmark", ctx, "user-1", "https://example.com", "demo bookmark", PrefixCreate).Return(&entity.Bookmark{
					Base:        base.Base{Id: "b-1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
					Url:         "https://example.com",
					Description: "demo bookmark",
					Code:        PrefixCreate + "ABC123XYZ9",
					UserId:      "user-1",
				}, nil)
				return repo
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
				assert.Contains(t, info.Code, PrefixCreate)
			},
		},
		{
			name: "repository error",
			setupRepo: func(ctx context.Context) *bookmark_mocks.BookmarkRepository {
				repo := bookmark_mocks.NewBookmarkRepository(t)
				repo.On("CreateBookmark", ctx, "user-1", "https://example.com", "demo bookmark", PrefixCreate).Return(nil, assert.AnError)
				return repo
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
			repo := tc.setupRepo(ctx)
			service := NewBookmarkService(repo)
			info, err := service.NewBookmark(ctx, "user-1", tc.input)
			tc.expectedFunc(t, info, err)
		})
	}
}
