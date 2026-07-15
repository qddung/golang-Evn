package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/homework/lab/internal/repository/mocks"
	mocksGenerateHelper "github.com/homework/lab/pkg/helpers/mocks"
)

var testErr = errors.New("test error")

const testExpTime = 60 * time.Second
const linkKeyLength = 7

func TestService_CreateShortenLink(t *testing.T) {
	testCases := []struct {
		name string

		setupRepo   func(ctx context.Context) *mocks.URLStorage
		setupKeyGen func() *mocksGenerateHelper.KeyGenerator

		expectedResult string
		expectedErr    error
	}{

		{
			name: "normal case - new key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "1234567", "https://google.com", testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("1234567")

				return mockKeyGen
			},

			expectedResult: "1234567",
			expectedErr:    nil,
		},
		{
			name: "normal case - random the same key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("https://example.com", nil)
				mock.On("GetURL", ctx, "2345678").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "2345678", "https://google.com", testExpTime).Return(nil)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("1234567").Once()
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("2345678").Once()

				return mockKeyGen
			},

			expectedResult: "2345678",
			expectedErr:    nil,
		},
		{
			name: "err case - can't put new key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", redis.Nil)
				mock.On("StoreURL", ctx, "1234567", "https://google.com", testExpTime).Return(testErr)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("1234567")

				return mockKeyGen
			},

			expectedResult: "",
			expectedErr:    testErr,
		},
		{
			name: "err case - can't get key",

			setupRepo: func(ctx context.Context) *mocks.URLStorage {
				mock := mocks.NewURLStorage(t)
				mock.On("GetURL", ctx, "1234567").Return("", testErr)

				return mock
			},
			setupKeyGen: func() *mocksGenerateHelper.KeyGenerator {
				mockKeyGen := mocksGenerateHelper.NewKeyGenerator(t)
				mockKeyGen.On("GenerateRandomCode", linkKeyLength).Return("1234567")

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

			testService := NewShorternUrl(mockRepo, keygenMockHelper)

			result, err := testService.ShortenUrlShortenUrl(ctx, "https://google.com", 60)
			assert.Equal(t, result, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}

}
