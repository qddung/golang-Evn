package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	"github.com/homework/lab/internal/models/dto/api"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
)

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
