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

func TestBookmarkEndpoint_CreateBookmark(t *testing.T) {
	testCases := []struct {
		name               string
		setupTestHTTP      func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedContain    string
	}{
		{
			name: "create bookmark successfully",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := general_helpers.MakeJSONRequest(http.MethodPost, buildBookmarkURL(""), bookmark_model.NewBookmarkRequest{
					Url:         "https://example.com/test-create",
					Description: "integration test create",
				})
				resp := httptest.NewRecorder()
				RouterSetupAuthorization(router, jwtMock, resp, req)
				return resp
			},
			expectedStatusCode: http.StatusOK,
			expectedContain:    "success create bookmark",
		},
		{
			name: "invalid url returns bad request",
			setupTestHTTP: func(router api.Engine, jwtMock *jwt_pkg.MockJwt) *httptest.ResponseRecorder {
				req := general_helpers.MakeJSONRequest(http.MethodPost, buildBookmarkURL(""), bookmark_model.NewBookmarkRequest{
					Url:         "not a url",
					Description: "bad request",
				})
				resp := httptest.NewRecorder()
				RouterSetupAuthorization(router, jwtMock, resp, req)
				return resp
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedContain:    "Url",
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
