package connection

import (
	"testing"

	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/homework/lab/pkg/sqldb"
)

func InitDBConnectorMock(t *testing.T) DBConnector {
	gormSql, err := sqldb.NewMiniPostgres(t)
	if err != nil {
		t.Fatal(err)
	}
	return &dbConnector{
		redisClient: redisPkg.InitMockRedis(t),
		db:          gormSql,
	}
}
