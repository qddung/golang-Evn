package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/models/dto/api"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	"github.com/homework/lab/pkg/response"
)

type UserHandler interface {
	Register(c *gin.Context)
}

type userHandler struct {
	svc user_service.UserService
}

func NewUserHandler(svc user_service.UserService) UserHandler {
	return &userHandler{svc: svc}
}

var createSuccessMessage = "User created successfully!"

// Register      Register link
// @Summary      Register user
// @Description  Register user
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        input body userModel.UserRegister true "Input required"
// @Success      200 {object} api.Response[userModel.UserInfo]
// @Router       /v1/users/register [post]

func (u *userHandler) Register(c *gin.Context) {
	userRequest := &userModel.UserRegister{}
	if err := c.ShouldBindJSON(userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": response.ToDataResponse[any](err)})
		return
	}
	createdUser, err := u.svc.Register(c, *userRequest)
	if err == user_service.EmailExistError || err == user_service.UserNameExistError {
		c.JSON(http.StatusConflict, gin.H{"error": response.ToDataResponse[any](err)})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": response.ToDataResponse[any](err)})
		return
	}

	res := &api.Response[userModel.UserInfo]{
		Message: createSuccessMessage,
		Data:    createdUser,
	}
	c.JSON(http.StatusOK, res)
}
