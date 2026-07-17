package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/service/mocks"
	"github.com/homework/lab/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_ShortenUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		inputBody             interface{}
		setupMock             func(mockSvc *mocks.ShorternUrl, ctx *gin.Context)
		expectedStatusCode    int
		expectedResult        interface{}
		expectedErrorContains string
	}{
		{
			name: "Case 1: Rút gọn URL thành công (200 OK)",
			inputBody: map[string]interface{}{
				"url": "https://www.google.com",
				"exp": 3600,
			},
			setupMock: func(mockSvc *mocks.ShorternUrl, ctx *gin.Context) {
				mockSvc.On("ShortenUrlShortenUrl", ctx, "https://www.google.com", int64(3600)).Return("abcde1", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResult: shortenResMessage{
				Message: "Shorten URL generated successfully!",
				Code:    "abcde1",
			},
		},
		{
			name: "Case 2: Lỗi validation - URL không hợp lệ (400 Bad Request)",
			inputBody: map[string]interface{}{
				"url": "invalid-url",
				"exp": 3600,
			},
			setupMock:             func(mockSvc *mocks.ShorternUrl, ctx *gin.Context) {},
			expectedStatusCode:    http.StatusBadRequest,
			expectedErrorContains: "Field validation for 'Url' failed on the 'url' tag",
		},
		{
			name: "Case 3: Lỗi validation - Thiếu tham số exp (400 Bad Request)",
			inputBody: map[string]interface{}{
				"url": "https://www.google.com",
			},
			setupMock:             func(mockSvc *mocks.ShorternUrl, ctx *gin.Context) {},
			expectedStatusCode:    http.StatusBadRequest,
			expectedErrorContains: "Field validation for 'Exp' failed on the 'required' tag",
		},
		{
			name: "Case 4: Lỗi validation - exp vượt quá giới hạn 604800 (400 Bad Request)",
			inputBody: map[string]interface{}{
				"url": "https://www.google.com",
				"exp": 999999,
			},
			setupMock:             func(mockSvc *mocks.ShorternUrl, ctx *gin.Context) {},
			expectedStatusCode:    http.StatusBadRequest,
			expectedErrorContains: "Field validation for 'Exp' failed on the 'lte' tag",
		},
		{
			name:                  "Case 5: Lỗi validation - JSON body sai định dạng (400 Bad Request)",
			inputBody:             "{invalid-json}",
			setupMock:             func(mockSvc *mocks.ShorternUrl, ctx *gin.Context) {},
			expectedStatusCode:    http.StatusBadRequest,
			expectedErrorContains: "invalid character",
		},
		{
			name: "Case 6: Lỗi từ tầng Service (500 Internal Server Error)",
			inputBody: map[string]interface{}{
				"url": "https://www.google.com",
				"exp": 3600,
			},
			setupMock: func(mockSvc *mocks.ShorternUrl, ctx *gin.Context) {
				mockSvc.On("ShortenUrlShortenUrl", ctx, "https://www.google.com", int64(3600)).Return("", errors.New("redis internal error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResult:     response.InternalErrResponse,
		},
	}

	for _, tc := range testCases {
		// Preserve the test context for the parallel loop
		tc := tc

		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()

			// Initialize a mocked Gin HTTP context
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			var bodyBytes []byte
			if strBody, ok := tc.inputBody.(string); ok {
				bodyBytes = []byte(strBody)
			} else {
				bodyBytes, _ = json.Marshal(tc.inputBody)
			}
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBuffer(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			// Create a mock service
			mockSvc := mocks.NewShorternUrl(testItem) // Use testItem for the correct context

			// Configure the mock behavior
			if tc.setupMock != nil {

				tc.setupMock(mockSvc, c)
			}

			// Initialize the handler and inject the mock service (DI)
			shortenHandler := NewShortenURL(mockSvc)

			// Execute the function under test
			shortenHandler.ShortenUrl(c)

			// Check the HTTP status code
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code)

			// Compare the returned JSON response accurately
			if tc.expectedStatusCode == http.StatusOK {
				var actualResponse shortenResMessage
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				assert.NoError(testItem, err)
				assert.Equal(testItem, tc.expectedResult, actualResponse)
			} else if tc.expectedStatusCode == http.StatusInternalServerError {
				var actualResponse response.Message
				err := json.Unmarshal(rec.Body.Bytes(), &actualResponse)
				assert.NoError(testItem, err)
				assert.Equal(testItem, tc.expectedResult, actualResponse)
			} else {
				assert.Contains(testItem, rec.Body.String(), tc.expectedErrorContains)
			}
		})
	}
}

func TestShortenURL_Redirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.ShorternUrl

		expectedStatus int
	}{
		{
			name: "normal case - success",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/link/redirect/1234567",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShorternUrl {
				serviceMock := mocks.NewShorternUrl(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("https://google.com", nil)
				return serviceMock
			},

			expectedStatus: http.StatusFound,
		},
		{
			name: "err case - wrong code",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/link/redirect",
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *mocks.ShorternUrl {
				serviceMock := mocks.NewShorternUrl(t)
				return serviceMock
			},

			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "err case - service err",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(
					http.MethodGet,
					"/link/redirect/1234567",
					nil,
				)
				ctx.Params = gin.Params{{Key: "code", Value: "1234567"}}
			},
			setupMockService: func(ctx context.Context) *mocks.ShorternUrl {
				serviceMock := mocks.NewShorternUrl(t)
				serviceMock.On("GetLinkFromCode", ctx, "1234567").Return("", errors.New("test err"))
				return serviceMock
			},

			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockService := tc.setupMockService(gc)
			handler := NewShortenURL(mockService)
			handler.Redirect(gc)
			assert.Equal(t, tc.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
