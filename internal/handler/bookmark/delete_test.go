package bookmark_handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	bookmark_cache_mocks "github.com/homework/lab/internal/cache/bookmark/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_DeleteBookmark(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(ctx *gin.Context) *bookmark_cache_mocks.BookmarkCache
		withClaims   bool
		expectedCode int
		expectedText string
	}{
		{
			name: "success",
			setupMock: func(ctx *gin.Context) *bookmark_cache_mocks.BookmarkCache {
				mockSvc := bookmark_cache_mocks.NewBookmarkCache(t)
				mockSvc.On("DeleteBookmark", ctx, "user-1", "bookmark-1").Return(nil)
				return mockSvc
			},
			withClaims:   true,
			expectedCode: http.StatusOK,
			expectedText: "Delete bookmark successfully",
		},
		{
			name: "missing claims",
			setupMock: func(ctx *gin.Context) *bookmark_cache_mocks.BookmarkCache {
				return bookmark_cache_mocks.NewBookmarkCache(t)
			},
			withClaims:   false,
			expectedCode: http.StatusUnauthorized,
			expectedText: "Claims extract token",
		},
		{
			name: "service error",
			setupMock: func(ctx *gin.Context) *bookmark_cache_mocks.BookmarkCache {
				mockSvc := bookmark_cache_mocks.NewBookmarkCache(t)
				mockSvc.On("DeleteBookmark", ctx, "user-1", "bookmark-1").Return(response.NotFoundError)
				return mockSvc
			},
			withClaims:   true,
			expectedCode: http.StatusNotFound,
			expectedText: response.InternalError.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/bookmark-1", nil)
			ctx.Params = gin.Params{{Key: "id", Value: "bookmark-1"}}
			if tc.withClaims {
				ctx.Set("claims", jwt.MapClaims{"sub": "user-1"})
			}
			service := tc.setupMock(ctx)
			handler := NewBookmarkHandler(service)
			handler.DeleteBookmark(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedText)
		})
	}
}
