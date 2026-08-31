package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	"github.com/homework/lab/pkg/request_ultls"
)

// LoginLink Login
// @Summary Login
// @Description Login
// @Tags user
// @Accept json
// @Produce json
// @Param input body userModel.UserLogin true "Input required"
// @Success 200 {object} loginResponse
// @Failure      400  {object}  loginResponse
// @Failure      500  {object}  loginResponse
// @Router /v1/users/login [post]
func (u *userHandler) Login(c *gin.Context) {
	userRequest, err := request_ultls.ModelBindValidation[userModel.UserLogin](c)
	loginResponse := &loginResponse{}
	if err != nil {
		loginResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, loginResponse)
		return
	}

	token, err := u.svc.Login(c, *userRequest)

	if err != nil {
		if user_service.CheckErrorIsServiceErr(err) {
			loginResponse.Message = err.Error()
			c.JSON(http.StatusBadRequest, loginResponse)
			return
		}
		loginResponse.Message = "Internal server error"
		c.JSON(http.StatusInternalServerError, loginResponse)
		return
	}
	loginResponse.Token = token
	loginResponse.Message = "Login successfully"
	c.JSON(http.StatusOK, loginResponse)
}

type loginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
