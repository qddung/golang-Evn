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
		setup         func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService, userId, bookmarkId string)
		expectedError error
	}{
		{
			name: "success",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService, userId, bookmarkId string) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey(userId)).Return(nil)
				serviceMock.On("DeleteBookmark", ctx, userId, bookmarkId).Return(nil)
			},
		},
		{
			name: "connection error",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService, userId, bookmarkId string) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey(userId)).Return(redis.ErrClosed)
			},
			expectedError: redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userId, bookmarkId := "user-1", "bookmark-1"
			cacheMock := mocks.NewCache(t)
			serviceMock := bookmark_service_mocks.NewBookmarkService(t)
			tc.setup(cacheMock, serviceMock, userId, bookmarkId)
			instance := NewBookmarkCacheInstance(serviceMock, cacheMock)
			err := instance.DeleteBookmark(ctx, userId, bookmarkId)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
