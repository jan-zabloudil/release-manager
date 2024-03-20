package utils

import (
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

	registerCustomValidations()
}

func TranslateValidationErrs(errs validator.ValidationErrors) map[string]string {
	m := make(map[string]string)

	for _, err := range errs {
		field := strings.ToLower(err.Field())
		m[field] = err.Translate(translator)
	}
	return m
}

func registerTranslation(tag, msg string) {
	if err := Validate.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
		return ut.Add(tag, msg, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	}); err != nil {
		panic(err)
	}
}
