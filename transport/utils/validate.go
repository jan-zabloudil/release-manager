package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate   *validator.Validate
	translator ut.Translator
)

func init() {
	Validate = validator.New()
	english := en.New()
	uni := ut.New(english, english)
	translator, _ = uni.GetTranslator("en")
	_ = entrans.RegisterDefaultTranslations(Validate, translator)

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func TranslateValidationErrs(errs validator.ValidationErrors) map[string]string {
	m := make(map[string]string)

	for _, err := range errs {
		field := strings.ToLower(err.Field())
		m[field] = err.Translate(translator)
	}
	return m
}
