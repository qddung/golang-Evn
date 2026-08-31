package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	"github.com/homework/lab/internal/models/dto/api"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
)

// GetUserInfo
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} api.Response[userModel.UserInfo]
// @Failure      400  {object}  api.Response[any]
// @Failure      500  {object}  api.Response[any]
// @Router /v1/users [get]
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func (u *userHandler) GetUserInfo(c *gin.Context) {

	claims, err := authorization.GetClaims(c)

	apiResponse := &api.Response[userModel.UserInfo]{}
	if err != nil {
		apiResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, apiResponse)
		return
	}

	userInfo, err := u.svc.GetUserInfo(c, claims["sub"].(string))
	if user_service.CheckErrorIsServiceErr(err) {
		apiResponse.Message = err.Error()
		c.JSON(http.StatusBadRequest, apiResponse)
		return
	}
	if err != nil {
		apiResponse.Message = "Internal server error"
		c.JSON(http.StatusInternalServerError, apiResponse)
		return
	}
	apiResponse.Data = userInfo
	c.JSON(http.StatusOK, apiResponse)
}
