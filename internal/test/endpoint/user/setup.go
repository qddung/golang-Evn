package user_endpoint

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	"github.com/homework/lab/internal/test/data/fixture"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
)

func BuildApiEngine(t *testing.T, fix fixture.Fixture, engOpt *api.EnginOpt) api.Engine {
	connectorMock, errConnector := connection.InitDBConnectorMock(t, fix)
	if errConnector != nil {
		t.Fatal(errConnector)
	}
	engOpt.App = gin.New()
	engOpt.Connector = connectorMock
	apiEngine := api.NewEngine(engOpt)
	return apiEngine
}

func BuildUserHandlerFull(testItem *testing.T, cfg *config.Config) api.Engine {
	db := fixture.NewUserTestCase(testItem)
	jwtMock := jwt_pkg.NewMockJwt()
	jwtGenerator := jwtMock.JwtGenarate
	jwtValidator := jwtMock.JwtValidate
	apiEngine := BuildApiEngine(testItem, db, &api.EnginOpt{
		Cfg:          cfg,
		JwtGenerator: jwtGenerator,
		JwtValidator: jwtValidator,
	})
	return apiEngine

}
