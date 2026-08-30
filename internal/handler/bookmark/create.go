package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homework/lab/internal/handler/authorization"
	"github.com/homework/lab/internal/models/dto/api"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/request_ultils"
	"github.com/homework/lab/pkg/response"
)

var SuccessCreateBookmark = "success create bookmark"

func (h *bookmarkHandler) CreateBookmark(c *gin.Context) {
	request, err := request_ultils.ModelBindValidation[bookmark_model.NewBookmarkRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ToDataResponse[*bookmark_model.BookmarkInfo](err))
		return
	}

	claims, err := authorization.GetClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ToDataResponse[*bookmark_model.BookmarkInfo](authorization.ClaimsNotFound))
		return
	}
	userId := claims["sub"].(string)

	bookmark, err := h.svc.NewBookmark(c, userId, request)

	res := &api.Response[bookmark_model.BookmarkInfo]{}
	if err != nil {
		res.Message = err.Error()
		c.JSON(response.MapErrorToHttpCode[err], res)
		return
	}
	res.Data = bookmark
	res.Message = SuccessCreateBookmark
	c.JSON(http.StatusOK, res)
}
