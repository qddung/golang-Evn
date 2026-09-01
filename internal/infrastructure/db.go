package infrastructure

import (
	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/homework/lab/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateSqlClient() *gorm.DB {
	sqlClient, err := sqldb.NewSqlDB()
	if err != nil {
		panic(err)
	}
	return sqlClient
}
func CreateRedisClient() *redis.Client {
	rdClient, err := redisPkg.NewRedisClient()
	if err != nil {
		panic(err)
	}
	return rdClient
}
