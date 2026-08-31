package bookmark_endpoint

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	"github.com/homework/lab/internal/test/data/fixture"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
)

const testBookmarkUserID = "4e90220a-51f6-49e4-bc0e-44e2f321475a"

func buildBookmarkIntegrationConfig() *config.Config {
	return &config.Config{
		AppPort:     "8080",
		ServiceName: "app_service",
		InstanceID:  "instance_01",
	}
}

func BuildBookmarkApiEngine(t *testing.T, fix fixture.Fixture, jwtMock *jwt_pkg.MockJwt) api.Engine {
	t.Helper()
	connectorMock, err := connection.InitDBConnectorMock(t, fix)
	if err != nil {
		t.Fatal(err)
	}

	engine := api.NewEngine(&api.EnginOpt{
		App:          gin.New(),
		Cfg:          buildBookmarkIntegrationConfig(),
		Connector:    connectorMock,
		JwtGenerator: jwtMock.JwtGenarate,
		JwtValidator: jwtMock.JwtValidate,
	})

	return engine
}

var dummyUserId = "4e90220a-51f6-49e4-bc0e-44e2f321475a"

func RouterSetupAuthorization(eng api.Engine, jwtMock *jwt_pkg.MockJwt, w http.ResponseWriter, req *http.Request) {
	claims := jwt.MapClaims{
		"sub": dummyUserId,
	}
	token, err := jwtMock.JwtGenarate.GenerateToken(claims)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	eng.ServeHTTP(w, req)
}

func buildBookmarkURL(path string) string {
	return fmt.Sprintf("/v1/bookmarks%s", path)
}
