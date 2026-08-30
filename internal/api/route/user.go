package route

import (
	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/api/middleware"
	user_handler "github.com/homework/lab/internal/handler/user"
)

func UserRoute(api *gin.RouterGroup, handler user_handler.UserHandler, jwtMiddeleware *middleware.JwtAuthMiddleware) {
	api.POST("/users/login", handler.Login)
	api.POST("/users/register", handler.Register)

	api.Use(jwtMiddeleware.JwtAuth()) // middelware
	api.GET("/self/info", handler.GetUserInfo)
	api.PUT("/self/info", handler.UpdateUserInfo)
}
