package bookmark_endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/homework/lab/internal/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/internal/test/data/fixture"
	general_helpers "github.com/homework/lab/internal/test/general"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkEndpoint_UpdateBookmark(t *testing.T) {
	testCases := []struct {
		name               string
		setupTestHTTP      func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedContain    string
	}{
		{
			name: "update bookmark successfully",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := general_helpers.MakeJSONRequest(http.MethodPut, buildBookmarkURL("/28d51812-2c60-4896-853b-cfccd06d5243"), bookmark_model.UpdateBookmarkRequest{
					Url:         "https://example.com/updated",
					Description: "updated description",
				})
				resp := httptest.NewRecorder()
				RouterSetupAuthorization(router, jwtMock, resp, req)
				return resp
			},
			expectedStatusCode: http.StatusOK,
			expectedContain:    "Update bookmark successfully",
		},
		{
			name: "update unknown bookmark returns not found",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := general_helpers.MakeJSONRequest(http.MethodPut, buildBookmarkURL("/99999999-9999-9999-9999-999999999999"), bookmark_model.UpdateBookmarkRequest{
					Url:         "https://example.com/updated",
					Description: "updated description",
				})
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
