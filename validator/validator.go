package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Validate = validator.Validate

func ValidateStruct(c *gin.Context, req any, v *Validate) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		return false
	}

	if err := v.Struct(req); err != nil {
		return false
	}

	return true
}

func ValidateStructOnly(req any, v *Validate) error {
	return v.Struct(req)
}
