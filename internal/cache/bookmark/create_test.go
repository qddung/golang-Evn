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
		setup         func(*mocks.Cache, *bookmark_service_mocks.BookmarkService)
		expected      *bookmark_model.BookmarkInfo
		expectedError error
	}{
		{
			name: "success",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey("user-1")).Return(nil)
				serviceMock.On("NewBookmark", ctx, "user-1", request).Return(expected, nil)
			},
			expected: expected,
		},
		{
			name: "connection error",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey("user-1")).Return(redis.ErrClosed)
			},
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cacheMock := mocks.NewCache(t)
			serviceMock := bookmark_service_mocks.NewBookmarkService(t)
			tc.setup(cacheMock, serviceMock)
			instance := &BookmarkCacheInstance{service: serviceMock, cache: cacheMock}

			result, err := instance.NewBookmark(ctx, "user-1", request)

			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, result)
		})
	}
}
