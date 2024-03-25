package utils

import (
	svcmodel "release-manager/service/model"
	"release-manager/transport/model"

	"github.com/go-playground/validator/v10"
)

func ValidateURL(fl validator.FieldLevel) bool {
	_, err := svcmodel.NewEnvURL(fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

func ValidateSourceCodeRequired(fl validator.FieldLevel) bool {
	sourceCode, ok := fl.Field().Interface().(model.SourceCode)
	if !ok {
		return false
	}

	if sourceCode.Tag == nil || sourceCode.TargetCommitIsh == nil {
		return false
	}

	return true
}

func ValidateSourceCodeIfPresent(fl validator.FieldLevel) bool {
	sourceCode, ok := fl.Field().Interface().(model.SourceCode)
	if !ok {
		return false
	}

	empty := model.SourceCode{}
	if sourceCode != empty {
		if sourceCode.Tag == nil || sourceCode.TargetCommitIsh == nil {
			return false
		}
	}

	return true
}

func registerCustomValidations() {
	envUrlTag := "env_url"
	sourceCodeRequired := "source_code_required"
	sourceCodeIfPresent := "source_code_if_present"

	if err := Validate.RegisterValidation(envUrlTag, ValidateURL); err != nil {
		panic(err)
	}
	registerTranslation(envUrlTag, "{0} must be valid absolute url")

	if err := Validate.RegisterValidation(sourceCodeRequired, ValidateSourceCodeRequired); err != nil {
		panic(err)
	}
	registerTranslation(sourceCodeRequired, "{0} is required field and must have 'tag' and 'target_commitish' fields")

	if err := Validate.RegisterValidation(sourceCodeIfPresent, ValidateSourceCodeIfPresent); err != nil {
		panic(err)
	}
	registerTranslation(sourceCodeIfPresent, "{0} is required field and must have 'tag' and 'target_commitish' fields")
}
