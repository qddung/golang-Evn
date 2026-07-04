package main

import (
	"github.com/homework/lab/internal/api"
	"github.com/homework/lab/internal/config"
)

func main() {
	// get app config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	apiEngine := api.NewEngine(cfg)
	err = apiEngine.Run()
	if err != nil {
		panic(err)
	}
}
