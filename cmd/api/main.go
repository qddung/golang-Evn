package main

import (
	"github.com/homework/lab/internal/infrastructure"
)

// @title Book API
// @version 1.0
// @description This is a book API
// @BasePath /
// initRoutes initializes the routes for the app engine
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
func main() {
	apiEngine := infrastructure.CreateApi()
	err := apiEngine.Run()
	if err != nil {
		panic(err)
	}
}
