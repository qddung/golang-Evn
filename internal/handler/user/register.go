package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/models/dto/api"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	"github.com/homework/lab/pkg/response"
)

// RegisterLink     godoc
// @Summary      Register user
// @Description  Register user
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        input body userModel.UserRegister true "Input required"
// @Success      200 {object} api.Response[userModel.UserInfo]
// @Failure      400  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router       /v1/users/register [post]
func (u *userHandler) Register(c *gin.Context) {
	userRequest := &userModel.UserRegister{}
	if err := c.ShouldBindJSON(userRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.ToDataResponse[any](err))
		return
	}
	createdUser, err := u.svc.Register(c, *userRequest)
	if err == user_service.ServiceErr.EmailExistError || err == user_service.ServiceErr.UserNameExistError {
		c.JSON(http.StatusConflict, response.ToDataResponse[any](err))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ToDataResponse[any](err))
		return
	}

	res := &api.Response[userModel.UserInfo]{
		Message: createSuccessMessage,
		Data:    createdUser,
	}
	c.JSON(http.StatusOK, res)
}
