package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateStruct(c *gin.Context, req any, v *validator.Validate) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		return false
	}

	if err := v.Struct(req); err != nil {
		return false
	}

	return true
}

func ValidateStructOnly(req any, v *validator.Validate) error {
	return v.Struct(req)
}
