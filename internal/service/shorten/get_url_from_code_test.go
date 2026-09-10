package shorten_service

import (
	"context"
	"testing"

	"github.com/homework/lab/internal/models/entity"
	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	"github.com/homework/lab/internal/repository/shorten/mocks"
	bookmark_service "github.com/homework/lab/internal/service/bookmark"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestService_GetLinkFromCode_CallShortenCode(t *testing.T) {
	testCases := []struct {
		name           string
		code           string
		setupRepo      func(ctx context.Context, code string) *mocks.URLStorage
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case",
			code: prefix + "123456",
			setupRepo: func(ctx context.Context, code string) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, code).Return("https://google.com", nil)

				return mock
			},
			expectedResult: "https://google.com",
			expectedErr:    nil,
		},
		{
			name: "err case - can't get key",
			code: prefix + "123456",
			setupRepo: func(ctx context.Context, code string) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, code).Return("", testErr)
				return mock
			},
			expectedResult: "",
			expectedErr:    testErr,
		},

		// {
		// 	name: "err case - can't get key with prefix",
		// 	code: "123456",
		// 	setupRepo: func(ctx context.Context, code string) *mocks.URLStorage {
		// 		return nil
		// 	},
		// 	expectedResult: "",
		// 	expectedErr:    testErr,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx, tc.code)

			testService := NewShorternUrl(mockRepo, nil, nil)

			result, err := testService.GetLinkFromCode(ctx, tc.code)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}

func TestService_GetLinkFromCode_CallBookmarkCode(t *testing.T) {
	testCases := []struct {
		name           string
		code           string
		setupRepo      func(ctx context.Context, code string) *bookmark_mocks.BookmarkRepository
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case with prefix bookmark",
			code: bookmark_service.PrefixCreate + "1234567890",
			setupRepo: func(ctx context.Context, code string) *bookmark_mocks.BookmarkRepository {
				mock := bookmark_mocks.NewBookmarkRepository(t)
				bookmark_return := &entity.Bookmark{
					Code: code,
					Url:  "https://google.com",
				}
				mock.On("FindBookmarkByCode", ctx, code).Return(bookmark_return, nil)
				return mock
			},
			expectedResult: "https://google.com",
			expectedErr:    nil,
		},
		{
			name: "err case - can't get key in repository",
			code: bookmark_service.PrefixCreate + "1234567890",
			setupRepo: func(ctx context.Context, code string) *bookmark_mocks.BookmarkRepository {
				mock := bookmark_mocks.NewBookmarkRepository(t)
				mock.On("FindBookmarkByCode", ctx, code).Return(nil, redis.Nil)
				return mock
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
		{
			name: "err case - can't get key by prefix",
			code: "1234567890",
			setupRepo: func(ctx context.Context, code string) *bookmark_mocks.BookmarkRepository {
				return nil
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			bookmarkMockRepo := tc.setupRepo(ctx, tc.code)
			testService := NewShorternUrl(nil, nil, bookmarkMockRepo)

			result, err := testService.GetLinkFromCode(ctx, tc.code)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}
