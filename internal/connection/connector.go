package connection

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBConnector interface {
	GetRedisClient() *redis.Client
	GetSqlDB() *gorm.DB
}

type dbConnector struct {
	redisClient *redis.Client
	db          *gorm.DB
}

func NewDBConnector(redisClient *redis.Client, db *gorm.DB) DBConnector {
	return &dbConnector{
		redisClient: redisClient,
		db:          db,
	}
}

func (db *dbConnector) GetRedisClient() *redis.Client {
	return db.redisClient
}

func (db *dbConnector) GetSqlDB() *gorm.DB {
	return db.db
}
