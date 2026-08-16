package user_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userModel "github.com/homework/lab/internal/models/dto/api/user"
	user_service "github.com/homework/lab/internal/service/user"
	"github.com/homework/lab/pkg/request_ultls"
)

func (u *userHandler) Login(c *gin.Context) {
	userRequest, err := request_ultls.ModelBindValidation[userModel.UserLogin](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := u.svc.Login(c, *userRequest)

	if err != nil {
		if user_service.CheckErrorIsServiceErr(err) {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
