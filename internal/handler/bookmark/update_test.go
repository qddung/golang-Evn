package bookmark_handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	bookmark_service_mocks "github.com/homework/lab/internal/service/bookmark/mocks"
	"github.com/homework/lab/pkg/request_ultils"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	request_ultils.InputValidator.RegisterValidation("url", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok || value == "" {
			return false
		}
		_, err := url.ParseRequestURI(value)
		return err == nil
	})
	request_ultils.InputValidator.RegisterValidation("url;required", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok || value == "" {
			return false
		}
		_, err := url.ParseRequestURI(value)
		return err == nil
	})
}

func TestBookmarkHandler_UpdateBookmark(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService
		body         string
		withClaims   bool
		expectedCode int
		expectedText string
	}{
		{
			name: "success",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				mockSvc := bookmark_service_mocks.NewBookmarkService(t)
				mockSvc.On("UpdateBookmark", ctx, mock.MatchedBy(func(req *bookmark_model.UpdateBookmarkRequest) bool {
					return req != nil && req.Url == "https://example.com/updated" && req.Description == "updated demo"
				}), "user-1", "bookmark-1").Return(nil)
				return mockSvc
			},
			body:         `{"url":"https://example.com/updated","description":"updated demo"}`,
			withClaims:   true,
			expectedCode: http.StatusOK,
			expectedText: "Update bookmark successfully",
		},
		{
			name: "missing claims",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				return bookmark_service_mocks.NewBookmarkService(t)
			},
			body:         `{"url":"https://example.com/updated","description":"updated demo"}`,
			withClaims:   false,
			expectedCode: http.StatusUnauthorized,
			expectedText: "Claims extract token",
		},
		{
			name: "invalid request",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				return bookmark_service_mocks.NewBookmarkService(t)
			},
			body:         `{"url":"not a url"}`,
			withClaims:   true,
			expectedCode: http.StatusBadRequest,
			expectedText: "Field validation",
		},
		{
			name: "service error",
			setupMock: func(ctx *gin.Context) *bookmark_service_mocks.BookmarkService {
				mockSvc := bookmark_service_mocks.NewBookmarkService(t)
				mockSvc.On("UpdateBookmark", ctx, mock.MatchedBy(func(req *bookmark_model.UpdateBookmarkRequest) bool {
					return req != nil && req.Url == "https://example.com/updated" && req.Description == "updated demo"
				}), "user-1", "bookmark-1").Return(response.NotFoundError)
				return mockSvc
			},
			body:         `{"url":"https://example.com/updated","description":"updated demo"}`,
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
			ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/bookmark-1", bytes.NewBufferString(tc.body))
			ctx.Request.Header.Set("Content-Type", "application/json")
			ctx.Params = gin.Params{{Key: "id", Value: "bookmark-1"}}
			if tc.withClaims {
				ctx.Set("claims", jwt.MapClaims{"sub": "user-1"})
			}
			service := tc.setupMock(ctx)
			handler := NewBookmarkHandler(service)
			handler.UpdateBookmark(ctx)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedText)
		})
	}
}
