package validator

import "github.com/go-playground/validator/v10"

var instance *validator.Validate

func GetValidator() *validator.Validate {
	if instance == nil {
		instance = validator.New()
	}
	return instance
}

func InitValidator() error {
	instance = validator.New()
	return nil
}

func RegisterValidation(tag string, fn validator.Func) error {
	return GetValidator().RegisterValidation(tag, fn)
}
