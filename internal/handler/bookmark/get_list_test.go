package bookmark_handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	api "github.com/homework/lab/internal/models/dto/api"
	bookmark_service_mocks "github.com/homework/lab/internal/service/bookmark/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService
		query        string
		withClaims   bool
		expectedCode int
		expectedText string
	}{
		{
			name: "success",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				mockSvc := bookmark_service_mocks.NewBookmarkService(t)
				mockSvc.On("GetBookmarks", ctx, "user-1", mock.MatchedBy(func(query *bookmark_model.GetBookmarksQuery) bool {
					return query != nil && query.Page == 1 && query.Limit == 10 && query.Sort == "created_at desc"
				})).Return(&api.PaginatedResponse[bookmark_model.BookmarkInfo]{
					Data: []bookmark_model.BookmarkInfo{{Id: "b-1", Url: "https://example.com", Description: "demo", Code: "ABC123"}},
					Pagination: api.Pagination{Page: 1, Limit: 10, Total: 1},
				}, nil)
				return mockSvc
			},
			query:        "?page=1&limit=10&sort=created_at+desc",
			withClaims:   true,
			expectedCode: http.StatusOK,
			expectedText: "b-1",
		},
		{
			name: "missing claims",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				return bookmark_service_mocks.NewBookmarkService(t)
			},
			query:        "?page=1&limit=10",
			withClaims:   false,
			expectedCode: http.StatusUnauthorized,
			expectedText: "Claims extract token",
		},
		{
			name: "service error",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				mockSvc := bookmark_service_mocks.NewBookmarkService(t)
				mockSvc.On("GetBookmarks", ctx, "user-1", mock.MatchedBy(func(query *bookmark_model.GetBookmarksQuery) bool {
					return query != nil && query.Page == 1 && query.Limit == 10 && query.Sort == "created_at desc"
				})).Return((*api.PaginatedResponse[bookmark_model.BookmarkInfo])(nil), response.InternalError)
				return mockSvc
			},
			query:        "?page=1&limit=10&sort=created_at+desc",
			withClaims:   true,
			expectedCode: http.StatusInternalServerError,
			expectedText: response.InternalError.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks"+tc.query, nil)
			if tc.withClaims {
				ctx.Set("claims", jwt.MapClaims{"sub": "user-1"})
			}
			service := tc.setupMock(ctx)
			handler := NewBookmarkHandler(service)
			handler.GetBookmarks(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedText)
		})
	}
}
