package utils

import (
	svcmodel "release-manager/service/model"

	"github.com/go-playground/validator/v10"
)

func ValidateURL(fl validator.FieldLevel) bool {
	_, err := svcmodel.NewEnvURL(fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

func registerCustomValidations() {
	envUrlTag := "env_url"

	if err := Validate.RegisterValidation(envUrlTag, ValidateURL); err != nil {
		panic(err)
	}
	registerTranslation(envUrlTag, "{0} must be valid absolute url")
}
