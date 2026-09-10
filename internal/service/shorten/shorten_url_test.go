package shorten_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/homework/lab/internal/repository/shorten/mocks"
	base62_helper_mocks "github.com/homework/lab/pkg/helpers/base62/mocks"
)

var testErr = errors.New("test error")

const testExpTime = 60 * time.Second
const linkKeyLength = 6

func TestService_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name        string
		setupRepo   func(ctx context.Context, url string) *mocks.URLStorage
		setupKeyGen func(url string) *base62_helper_mocks.Base62Helper

		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - new key",
			setupRepo: func(ctx context.Context, url string) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, prefix+"123456").Return("", redis.Nil)
				mock.On("StoreURL", ctx, prefix+"123456", url, testExpTime).Return(nil)
				return mock
			},
			setupKeyGen: func(url string) *base62_helper_mocks.Base62Helper {
				mockKeyGen := base62_helper_mocks.NewBase62Helper(t)
				mockKeyGen.On("Encode", url).Return("123456")
				return mockKeyGen
			},

			expectedResult: "123456",
			expectedErr:    nil,
		},
		{
			name: "normal case - random the same key",
			setupRepo: func(ctx context.Context, url string) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, prefix+"234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, prefix+"234567", url, testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func(url string) *base62_helper_mocks.Base62Helper {
				mockKeyGen := base62_helper_mocks.NewBase62Helper(t)
				mockKeyGen.On("Encode", url).Return("234567")
				return mockKeyGen
			},
			expectedResult: "234567",
			expectedErr:    nil,
		},
		{
			name: "err case - can't put new key",
			setupRepo: func(ctx context.Context, url string) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, prefix+"123456").Return("", redis.Nil)
				mock.On("StoreURL", ctx, prefix+"123456", url, testExpTime).Return(testErr)

				return mock
			},
			setupKeyGen: func(url string) *base62_helper_mocks.Base62Helper {
				mockKeyGen := base62_helper_mocks.NewBase62Helper(t)
				mockKeyGen.On("Encode", url).Return("123456")

				return mockKeyGen
			},

			expectedResult: "",
			expectedErr:    testErr,
		},
		{
			name: "err case - can't get key",
			setupRepo: func(ctx context.Context, url string) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, prefix+"123456").Return("", testErr)
				return mock
			},
			setupKeyGen: func(url string) *base62_helper_mocks.Base62Helper {
				mockKeyGen := base62_helper_mocks.NewBase62Helper(t)
				mockKeyGen.On("Encode", url).Return("123456")
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
			url := "https://google.com"

			mockRepo := tc.setupRepo(ctx, url)
			keygenMockHelper := tc.setupKeyGen(url)
			testService := NewShorternUrl(mockRepo, keygenMockHelper, nil)

			result, err := testService.ShortenUrlShortenUrl(ctx, url, 60)
			if tc.expectedResult != "" {
				assert.Equal(t, result, prefix+tc.expectedResult)
			}
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}
