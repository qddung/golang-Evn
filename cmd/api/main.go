package main

import (
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/homework/lab/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// @title Book API
// @version 1.0
// @description This is a book API
// @BasePath /
func main() {
	// get app config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// create redis client
	rdClient := createRedisClient()

	// create SQL client
	sqlClient := createSqlClient()

	// connector
	connector := connection.NewDBConnector(rdClient, sqlClient)

	apiEngine := api.NewEngine(cfg, connector)
	err = apiEngine.Run()
	if err != nil {
		panic(err)
	}
}
func createSqlClient() *gorm.DB {
	sqlClient, err := sqldb.NewSqlDB()
	if err != nil {
		panic(err)
	}
	return sqlClient
}
func createRedisClient() *redis.Client {
	rdClient, err := redisPkg.NewRedisClient()
	if err != nil {
		panic(err)
	}
	return rdClient
}
