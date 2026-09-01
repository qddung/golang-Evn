package main

import (
	"github.com/homework/lab/internal/infrastructure"
)

// @title Book API
// @version 1.0
// @description This is a book API
// @BasePath /
// initRoutes initializes the routes for the app engine
// @securityDefinitions.apikey JWT
// @in header
// @name JWT
func main() {
	apiEngine := infrastructure.CreateApi()
	err := apiEngine.Run()
	if err != nil {
		panic(err)
	}
}
