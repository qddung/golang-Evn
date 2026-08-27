package bookmark_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bookmark_model "github.com/homework/lab/internal/models/dto/api/bookmark"
	"github.com/homework/lab/pkg/request_ultils"
	"github.com/homework/lab/pkg/response"
)

func (h *bookmarkHandler) CreateBookmark(c *gin.Context) {
	request, err := request_ultils.ModelBindValidation[bookmark_model.NewBookmarkRequest](c)
	if err != nil {
		c.JSON(response.MapErrorToHttpCode[err], gin.H{"error": err.Error()})
		return
	}
	res := h.svc.NewBookmark(c, request)
	c.JSON(http.StatusOK, res)
}
