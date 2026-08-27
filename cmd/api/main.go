package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/constant"
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	"github.com/homework/lab/internal/connection"
	"github.com/homework/lab/internal/models/entity"
	jwt_pkg "github.com/homework/lab/pkg/jwt"
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

	if errMigrate := sqlClient.AutoMigrate(&entity.User{}); errMigrate != nil {
		log.Fatalf("Migration failed: %v", errMigrate)
	}

	// connector
	connector := connection.NewDBConnector(rdClient, sqlClient)
	jwtGenerator := jwt_pkg.NewJWTGenerator(constant.PrivateKeyPath)
	jwtValidator := jwt_pkg.NewJWTValidator(constant.PublicKeyPath)
	apiEngine := api.NewEngine(&api.EnginOpt{
		App:          gin.New(),
		Cfg:          cfg,
		Connector:    connector,
		JwtGenerator: jwtGenerator,
		JwtValidator: jwtValidator,
	})
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
