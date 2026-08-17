package user_endpoint

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	"github.com/homework/lab/internal/models/dto/api/user"
	general_helpers "github.com/homework/lab/internal/test/general"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestService_Login(t *testing.T) {
	testCases := []struct {
		name               string
		setupTestHTTP      func(router api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		configTest         *config.Config
	}{
		{
			name: "Login successfully",
			setupTestHTTP: func(router api.Engine) *httptest.ResponseRecorder {
				reqPost := general_helpers.MakeJSONRequest(http.MethodPost, "/v1/users/login", user.UserLogin{
					Password: "12345678",
					UserName: "acd3",
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
			fmt.Printf("Loaded config: %+v\n", tc.configTest)
			connectorMock, errConnector := connection.InitDBConnectorMock(testItem)
			if errConnector != nil {
				testItem.Fatal(errConnector)
			}
			jwtMock := jwt_pkg.NewMockJwt()
			jwtGenerator := jwtMock.JwtGenarate
			jwtValidator := jwtMock.JwtValidate
			apiEngine := api.NewEngine(&api.EnginOpt{
				App:          gin.New(),
				Cfg:          tc.configTest,
				Connector:    connectorMock,
				JwtGenerator: jwtGenerator,
				JwtValidator: jwtValidator,
			})
			rec := tc.setupTestHTTP(apiEngine)
			// Check status code
			// assert.Equal(testItem, "", rec.Body.String())
			assert.Equal(testItem, tc.expectedStatusCode, rec.Code, "Expected status code does not match actual status code")
		})
	}
}
