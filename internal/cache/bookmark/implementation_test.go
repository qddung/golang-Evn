package bookmark_cache

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/repository/cache/mocks"
	bookmark_service_mocks "github.com/homework/lab/internal/service/bookmark/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkCacheInstance_Get(t *testing.T) {
	ctx := context.Background()
	query := &bookmark_model.GetBookmarksQuery{Page: 1, Limit: 10}
	expected := &api.PaginatedResponse[bookmark_model.BookmarkInfo]{
		Data:       []bookmark_model.BookmarkInfo{{Id: "bookmark-1", Code: "code-1"}},
		Pagination: api.Pagination{Page: 1, Limit: 10, Total: 1},
	}
	expectedBytes, err := json.Marshal(expected)
	assert.NoError(t, err)

	testCases := []struct {
		name          string
		setup         func(*mocks.Cache, *bookmark_service_mocks.BookmarkService)
		expected      *api.PaginatedResponse[bookmark_model.BookmarkInfo]
		expectedError error
	}{
		{
			name: "success",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("Read", ctx, GetGroupKey("user-1"), "1_10").Return(string(expectedBytes), nil)
			},
			expected: expected,
		},
		{
			name: "connection error",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				var key = GetGroupKey("user-1")
				cacheMock.On("Read", ctx, key, "1_10").Return("", redis.ErrClosed)
				serviceMock.On("GetBookmarks", ctx, "user-1", query).Return(expected, nil)
				cacheMock.On("Write", ctx, key, "1_10", mock.AnythingOfType("[]uint8"), ttl).Return(redis.ErrClosed)
			},
			expectedError: ErrorWriteCache,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cacheMock := mocks.NewCache(t)
			serviceMock := bookmark_service_mocks.NewBookmarkService(t)
			tc.setup(cacheMock, serviceMock)
			instance := &BookmarkCacheInstance{service: serviceMock, cache: cacheMock}

			result, err := instance.GetBookmarks(ctx, "user-1", query)

			assert.ErrorIs(t, err, tc.expectedError)
			if tc.expectedError == nil {
				assert.Equal(t, tc.expected, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

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

func TestBookmarkCacheInstance_Update(t *testing.T) {
	ctx := context.Background()
	request := &bookmark_model.UpdateBookmarkRequest{Url: "https://example.com/updated"}

	testCases := []struct {
		name          string
		setup         func(*mocks.Cache, *bookmark_service_mocks.BookmarkService)
		expectedError error
	}{
		{
			name: "success",
			setup: func(cacheMock *mocks.Cache, serviceMock *bookmark_service_mocks.BookmarkService) {
				cacheMock.On("DeleteGroup", ctx, GetGroupKey("user-1")).Return(nil)
				serviceMock.On("UpdateBookmark", ctx, request, "user-1", "bookmark-1").Return(nil)
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

			err := instance.UpdateBookmark(ctx, request, "user-1", "bookmark-1")

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

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
