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

func TestBookmarkCacheInstance_Update(t *testing.T) {
	ctx := context.Background()
	request := &bookmark_model.UpdateBookmarkRequest{Url: "https://example.com/updated"}

	testCases := []struct {
		name          string
		setup         func(userId, bookmarkId string, cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService)
		expectedError error
	}{
		{
			name: "success",
			setup: func(userId, bookmarkId string, cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey(userId)).Return(nil)
				serviceMock.On("UpdateBookmark", ctx, request, userId, bookmarkId).Return(nil)
			},
		},
		{
			name: "connection error",
			setup: func(userId, bookmarkId string, cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey(userId)).Return(redis.ErrClosed)
			},
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userId := "user-1"
			bookmarkId := "bookmark-1"
			cacheMock := mocks.NewCache(t)
			serviceMock := bookmark_service_mocks.NewBookmarkService(t)
			tc.setup(userId, bookmarkId, cacheMock, serviceMock)
			instance := &BookmarkCacheInstance{service: serviceMock, cache: cacheMock}
			err := instance.UpdateBookmark(ctx, request, userId, bookmarkId)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
