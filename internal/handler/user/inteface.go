package user_handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	user_service "github.com/homework/lab/internal/service/user"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserInfo(c *gin.Context)
	UpdateUserInfo(c *gin.Context)
}

type userHandler struct {
	svc user_service.UserService
}

var (
	ErrorExtractClaims = errors.New("error extract token")
	ClaimsNotFound     = errors.New("Claims extract token")
)

func GetClaims(c *gin.Context) (jwt.MapClaims, error) {
	tokenClaims, exists := c.Get("claims")

	if !exists {
		return nil, ClaimsNotFound
	}

	if _, ok := tokenClaims.(jwt.MapClaims); !ok {
		return nil, ErrorExtractClaims
	}
	return tokenClaims.(jwt.MapClaims), nil
}

func NewUserHandler(svc user_service.UserService) UserHandler {
	return &userHandler{svc: svc}
}

var createSuccessMessage = "User created successfully!"
