package user_handler

import (
	"github.com/gin-gonic/gin"
	user_service "github.com/homework/lab/internal/service/user"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type userHandler struct {
	svc user_service.UserService
}

func NewUserHandler(svc user_service.UserService) UserHandler {
	return &userHandler{svc: svc}
}

var createSuccessMessage = "User created successfully!"
