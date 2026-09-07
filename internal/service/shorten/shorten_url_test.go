package shorten_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/homework/lab/internal/models/entity"
	bookmark_mocks "github.com/homework/lab/internal/repository/bookmark/mocks"
	"github.com/homework/lab/internal/repository/shorten/mocks"
	mocksGenerateHelper "github.com/homework/lab/pkg/helpers/mocks"
)

var testErr = errors.New("test error")

const testExpTime = 60 * time.Second
const linkKeyLength = 6

func TestService_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name        string
		setupRepo   func(ctx context.Context) *mocks.URLStorage
		setupKeyGen func() *mocksGenerateHelper.KeyGenerator

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - new key",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "123456").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "123456", "https://google.com", testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("123456")

				return mockKeyGen
			},

			expectedResult: "123456",
			expectedErr:    nil,
		},
		{
			name: "normal case - random the same key",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "123456").Return("https://example.com", nil)
				mock.On("GetURL", ctx, "234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "234567", "https://google.com", testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("123456").Once()
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("234567").Once()

				return mockKeyGen
			},
			expectedResult: "234567",
			expectedErr:    nil,
		},
		{
			name: "err case - can't put new key",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "123456").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "123456", "https://google.com", testExpTime).Return(testErr)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("123456")

				return mockKeyGen
			},

			expectedResult: "",
			expectedErr:    testErr,
		},
		{
			name: "err case - can't get key",
			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "123456").Return("", testErr)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("123456")

				return mockKeyGen
			},

			expectedResult: "",
			expectedErr:    testErr,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)
			keygenMockHelper := tc.setupKeyGen()

			testService := NewShorternUrl(mockRepo, keygenMockHelper, nil)

			result, err := testService.ShortenUrlShortenUrl(ctx, "https://google.com", 60)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}

func TestService_GetLinkFromCode_CallShortenCode(t *testing.T) {
	testCases := []struct {
		name           string
		setupRepo      func(ctx context.Context) *mocks.URLStorage
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "123456").Return("https://google.com", nil)

				return mock
			},
			expectedResult: "https://google.com",
			expectedErr:    nil,
		},
		{
			name: "err case - can't get key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "123456").Return("", testErr)
				return mock
			},
			expectedResult: "",
			expectedErr:    testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			mockRepo := tc.setupRepo(ctx)

			testService := NewShorternUrl(mockRepo, nil, nil)

			result, err := testService.GetLinkFromCode(ctx, "123456")
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}

func TestService_GetLinkFromCode_CallBookmarkCode(t *testing.T) {
	testCases := []struct {
		name           string
		code           string
		setupRepo      func(ctx context.Context) *bookmark_mocks.BookmarkRepository
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case",
			code: "1234567890",
			setupRepo: func(ctx context.Context) *bookmark_mocks.BookmarkRepository {
				mock := bookmark_mocks.NewBookmarkRepository(t)
				bookmark_return := &entity.Bookmark{
					Code: "1234567890",
					Url:  "https://google.com",
				}
				mock.On("FindBookmarkByCode", ctx, "1234567890").Return(bookmark_return, nil)
				return mock
			},
			expectedResult: "https://google.com",
			expectedErr:    nil,
		},
		{
			name: "err case - can't get key",
			code: "1234567890",
			setupRepo: func(ctx context.Context) *bookmark_mocks.BookmarkRepository {
				mock := bookmark_mocks.NewBookmarkRepository(t)
				mock.On("FindBookmarkByCode", ctx, "1234567890").Return(nil, redis.Nil)
				return mock
			},
			expectedResult: "",
			expectedErr:    ErrCodeDoesntExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			bookmarkMockRepo := tc.setupRepo(ctx)
			testService := NewShorternUrl(nil, nil, bookmarkMockRepo)

			result, err := testService.GetLinkFromCode(ctx, tc.code)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}
