package main

import (
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
	redisPkg "github.com/homework/lab/pkg/redis"
	"github.com/redis/go-redis/v9"
)

// @title Book API
// @version 1.0
// @description This is a book API
// host localhost:8080
// @BasePath /
func main() {
	// get app config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// create redis client
	rdClient := createRedisClient()

	apiEngine := api.NewEngine(cfg, rdClient)
	err = apiEngine.Run()
	if err != nil {
		panic(err)
	}
}

func createRedisClient() *redis.Client {
	rdClient, err := redisPkg.NewRedisClient()
	if err != nil {
		panic(err)
	}
	return rdClient
}
