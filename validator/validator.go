package validator

import "github.com/gin-gonic/gin"

func ValidateStruct(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		return false
	}

	if err := GetValidator().Struct(req); err != nil {
		return false
	}

	return true
}

func ValidateStructOnly(req any) error {
	return GetValidator().Struct(req)
}
