package utils

import (
	"github.com/go-playground/validator/v10"
)

func ValidateField(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return err
		}
	}

	return nil
}
