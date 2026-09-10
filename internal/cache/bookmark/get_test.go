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
		setup         func() (*mocks.Cache, *bookmark_service_mocks.BookmarkService)
		expected      *api.PaginatedResponse[bookmark_model.BookmarkInfo]
		expectedError error
	}{
		{
			name: "success",
			setup: func() (*mocks.Cache, *bookmark_service_mocks.BookmarkService) {
				cacheMock := mocks.NewCache(t)
				serviceMock := bookmark_service_mocks.NewBookmarkService(t)
				cacheMock.On("Read", ctx, GetGroupKey("user-1"), "1_10").Return(string(expectedBytes), nil)
				return cacheMock, serviceMock
			},
			expected: expected,
		},
		{
			name: "connection error",
			setup: func() (*mocks.Cache, *bookmark_service_mocks.BookmarkService) {
				var key = GetGroupKey("user-1")

				cacheMock := mocks.NewCache(t)
				cacheMock.On("Read", ctx, key, "1_10").Return("", redis.ErrClosed)
				cacheMock.On("Write", ctx, key, "1_10", mock.AnythingOfType("[]uint8"), ttl).Return(redis.ErrClosed)

				serviceMock := bookmark_service_mocks.NewBookmarkService(t)
				serviceMock.On("GetBookmarks", ctx, "user-1", query).Return(expected, nil)
				return cacheMock, serviceMock
			},
			expectedError: ErrorWriteCache,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cacheMock, serviceMock := tc.setup()
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
