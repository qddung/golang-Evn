package bookmark_endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/test/data/fixture"
	general_helpers "github.com/homework/lab/internal/test/general"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkEndpoint_GetBookmarks(t *testing.T) {
	testCases := []struct {
		name               string
		setupTestHTTP      func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedContain    string
	}{
		{
			name: "get bookmarks successfully",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := general_helpers.MakeJSONRequest(http.MethodGet, buildBookmarkURL("?page=1&limit=10&sort=created_at+desc"), nil)
				resp := httptest.NewRecorder()
				RouterSetupAuthorization(router, jwtMock, resp, req)
				return resp
			},
			expectedStatusCode: http.StatusOK,
			expectedContain:    "28d51812-2c60-4896-853b-cfccd06d5243",
		},
		{
			name: "get bookmarks for unknown user returns empty list",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, buildBookmarkURL("?page=1&limit=10"), nil)
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
