package request_ultls

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var InputValidator = validator.New(validator.WithRequiredStructEnabled())

// ModelBindValidation validate request in json, uri, query, header from gin
func ModelBindValidation[T any](c *gin.Context) (*T, error) {
	model := new(T)
	// serialize full model
	if err := c.ShouldBindUri(model); err != nil {
		return nil, err
	}

	if err := c.ShouldBindQuery(model); err != nil {
		return nil, err
	}
	if err := c.ShouldBindHeader(model); err != nil {
		return nil, err
	}

	if c.Request.Method == "POST" {
		err := c.ShouldBindJSON(model)
		if err != nil {
			return nil, err
		}
	}

	// validate on validator-go
	if err := InputValidator.Struct(model); err != nil {
		return nil, err
	}
	return model, nil
}
