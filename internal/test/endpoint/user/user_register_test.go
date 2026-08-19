package user_endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/models/dto/api/user"
	"github.com/homework/lab/internal/test/data/fixture"
	general_helpers "github.com/homework/lab/internal/test/general"
	"github.com/stretchr/testify/assert"
)

func TestService_Register(t *testing.T) {
	testCases := []struct {
		name               string
		setupTestHTTP      func(router api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		configTest         *config.Config
	}{
		{
			name: "Register successfully",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqPost := general_helpers.MakeJSONRequest(http.MethodPost, "/v1/users/register", user.UserRegister{
					DisplayName: "test",
					Email:       "test@example.com",
					Password:    "123131242",
					UserName:    "testuser",
				})
				respPost := httptest.NewRecorder()
				router.ServeHTTP(respPost, reqPost)
				return respPost
			},
			expectedStatusCode: http.StatusOK,
			configTest: &config.Config{
				AppPort:     "8080",
				ServiceName: "app_service",
				InstanceID:  "instance_01",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(testItem *testing.T) {
			testItem.Parallel()
			fix := fixture.NewUserTestCase(testItem)
			apiEngine := BuildApiEngine(testItem, fix,
				&api.EnginOpt{
					Cfg: tc.configTest,
				})
			rec := tc.setupTestHTTP(apiEngine)
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code, "Expected status code does not match actual status code")
		})
	}

}
