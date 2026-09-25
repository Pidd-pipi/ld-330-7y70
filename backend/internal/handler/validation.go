package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// bindJSON decodes JSON and runs validator/v10 on DTO struct tags.
func bindJSON(c *gin.Context, obj interface{}) bool {
	if e := c.ShouldBindJSON(obj); e != nil {
		Fail(c, e)
		return false
	}
	if e := validate.Struct(obj); e != nil {
		Fail(c, fmt.Errorf("validation failed: %s", e))
		return false
	}
	return true
}
