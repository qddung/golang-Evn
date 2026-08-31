package main

import (
	"github.com/homework/lab/internal/infrastructure"
)

// @title Book API
// @version 1.0
// @description This is a book API
// @BasePath /
func main() {
	apiEngine := infrastructure.CreateApi()
	err := apiEngine.Run()
	if err != nil {
		panic(err)
	}
}
