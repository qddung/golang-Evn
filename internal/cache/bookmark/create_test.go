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
	expected := &bookmark_model.BookmarkInfo{Id: "bookmark-1"}

	testCases := []struct {
		name          string
		setup         func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService, userId string, request *bookmark_model.NewBookmarkRequest)
		expected      *bookmark_model.BookmarkInfo
		expectedError error
	}{
		{
			name: "success",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService, userId string, request *bookmark_model.NewBookmarkRequest) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey(userId)).Return(nil)
				serviceMock.On("NewBookmark", ctx, userId, request).Return(expected, nil)
			},
			expected: expected,
		},
		{
			name: "connection error",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService, userId string, request *bookmark_model.NewBookmarkRequest) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey(userId)).Return(redis.ErrClosed)
			},
			expectedError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, userId := &bookmark_model.NewBookmarkRequest{Url: "https://example.com"}, "user-1"
			cacheMock := mocks.NewCache(t)
			serviceMock := bookmark_service_mocks.NewBookmarkService(t)
			tc.setup(cacheMock, serviceMock, userId, request)
			instance := &BookmarkCacheInstance{service: serviceMock, cache: cacheMock}
			result, err := instance.NewBookmark(ctx, userId, request)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, result)
		})
	}
}
