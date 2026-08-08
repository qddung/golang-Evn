package connection

import (
	"testing"

	"github.com/homework/lab/internal/models/entity"
	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/homework/lab/pkg/sqldb"
)

func InitDBConnectorMock(t *testing.T) DBConnector {
	gormSql, err := sqldb.NewMiniPostgres(t) // gorm need run all fixture
	if err != nil {
		t.Fatal(err)
	}

	errGorm := gormSql.AutoMigrate(entity.User{})
	if errGorm != nil {
		t.Fatal(errGorm)
	}
	return &dbConnector{
		redisClient: redisPkg.InitMockRedis(t),
		db:          gormSql,
	}
}
