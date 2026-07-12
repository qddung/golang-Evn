package main

import (
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
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

	apiEngine := api.NewEngine(cfg)
	err = apiEngine.Run()
	if err != nil {
		panic(err)
	}
}
