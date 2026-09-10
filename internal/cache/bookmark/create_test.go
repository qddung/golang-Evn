package bookmark_cache

import (
	"context"
	"testing"

	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/repository/cache/mocks"
	bookmark_service_mocks "github.com/homework/lab/internal/service/bookmark/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkCacheInstance_Create(t *testing.T) {
	ctx := context.Background()
	request := &bookmark_model.NewBookmarkRequest{Url: "https://example.com"}
	expected := &bookmark_model.BookmarkInfo{Id: "bookmark-1"}

	testCases := []struct {
		name          string
		setup         func() (*mocks.Cache, *bookmark_service_mocks.BookmarkService)
		expected      *bookmark_model.BookmarkInfo
		expectedError error
	}{
		{
			name: "success",
			setup: func() (*mocks.Cache, *bookmark_service_mocks.BookmarkService) {
				cacheMock := mocks.NewCache(t)
				serviceMock := bookmark_service_mocks.NewBookmarkService(t)
				cacheMock.On("DeleteGroup", ctx, GetGroupKey("user-1")).Return(nil)
				serviceMock.On("NewBookmark", ctx, "user-1", request).Return(expected, nil)
				return cacheMock, serviceMock
			},
			expected: expected,
		},
		{
			name: "connection error",
			setup: func() (*mocks.Cache, *bookmark_service_mocks.BookmarkService) {
				cacheMock := mocks.NewCache(t)
				serviceMock := bookmark_service_mocks.NewBookmarkService(t)
				cacheMock.On("DeleteGroup", ctx, GetGroupKey("user-1")).Return(redis.ErrClosed)
				return cacheMock, serviceMock
			},
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cacheMock, serviceMock := tc.setup()
			instance := &BookmarkCacheInstance{service: serviceMock, cache: cacheMock}
			result, err := instance.NewBookmark(ctx, "user-1", request)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, result)
		})
	}
}
