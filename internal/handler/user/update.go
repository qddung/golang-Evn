package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
)

type updateUserReposonse struct {
	Message string
}

// UpdateUserInfo godoc
// @Summary Edit current user
// @Description Edit current user
// @Tags user
// @Accept json
// @Produce json
// @Param input body userModel.UpdateUserInput true "Input required"
// @Success 200 {object} updateUserReposonse
// @Failure      400  {object}  updateUserReposonse
// @Failure      500  {object}  updateUserReposonse
// @Router /v1/users [put]
func (u *userHandler) UpdateUserInfo(c *gin.Context) {
	request := &userModel.UpdateUserInput{}
	c.ShouldBindJSON(request)
	response := &updateUserReposonse{Message: "Edit current user successfully"}
	claims, err := GetClaims(c)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err = u.svc.UpdateUserInfo(c, claims["sub"].(string), request)
	if err != nil {
		response.Message = err.Error()
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
