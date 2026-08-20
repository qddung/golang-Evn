package connection

import (
	"testing"

	"github.com/homework/lab/internal/test/data/fixture"
	redisPkg "github.com/homework/lab/pkg/redis"
)

func InitDBConnectorMock(t *testing.T, fix fixture.Fixture) (DBConnector, error) {
	// fix := fixture.NewUserTestCase(t)
	gormSql := fixture.NewFixture(t, fix)
	return NewDBConnector(redisPkg.InitMockRedis(t), gormSql), nil
}
