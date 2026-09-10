package bookmark_cache

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/repository/cache/mocks"
	bookmark_service_mocks "github.com/homework/lab/internal/service/bookmark/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkCacheInstance_DeleteBookmark(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name          string
		setup         func(*mocks.Cache, *bookmark_service_mocks.BookmarkService)
		expectedError error
	}{
		{
			name: "success",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey("user-1")).Return(nil)
				serviceMock.On("DeleteBookmark", ctx, "user-1", "bookmark-1").Return(nil)
			},
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

			err := instance.DeleteBookmark(ctx, "user-1", "bookmark-1")

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
