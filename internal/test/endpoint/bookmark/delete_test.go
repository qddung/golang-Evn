package bookmark_endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/test/data/fixture"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkEndpoint_DeleteBookmark(t *testing.T) {
	testCases := []struct {
		name               string
		setupTestHTTP      func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedContain    string
	}{
		{
			name: "delete bookmark successfully",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, buildBookmarkURL("/28d51812-2c60-4896-853b-cfccd06d5243"), nil)
				resp := httptest.NewRecorder()
				RouterSetupAuthorization(router, jwtMock, resp, req)
				return resp
			},
			expectedStatusCode: http.StatusOK,
			expectedContain:    "Delete bookmark successfully",
		},
		{
			name: "delete unknown bookmark returns not found",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, buildBookmarkURL("/99999999-9999-9999-9999-999999999999"), nil)
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				return resp
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedContain:    "Unauthorized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			jwtMock := jwt_pkg.NewMockJwt()
			router := BuildBookmarkApiEngine(testItem, fixture.NewBookmarkTestCase(testItem), jwtMock)
			rec := tc.setupTestHTTP(router, jwtMock)
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code)
			assert.Contains(testItem, rec.Body.String(), tc.expectedContain)
		})
	}
}
